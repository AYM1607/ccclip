package server

import (
	"crypto/ecdh"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/AYM1607/ccclip/internal/config"
	"github.com/AYM1607/ccclip/internal/db"
	"github.com/AYM1607/ccclip/pkg/api"
	"github.com/AYM1607/ccclip/pkg/crypto"
)

func New(addr string) *http.Server {
	h := newHttpHandler()

	return &http.Server{
		Addr:    addr,
		Handler: h,
	}
}

type controller struct {
	store     db.DB
	publicKey *ecdh.PublicKey
	// TODO: This should not stay in memory for a long time.
	// keeping it as part of the controller for testing purposes only.
	privateKey *ecdh.PrivateKey
}

func newHttpHandler() http.Handler {
	r := mux.NewRouter()

	pbk, err := crypto.LoadPublicKey(config.Default.PublicKeyPath)
	if err != nil {
		panic("could not load server's public key")
	}
	pvk, err := crypto.LoadPrivateKey(config.Default.PrivateKeyPath)
	if err != nil {
		panic("could not load server's private key")
	}
	c := &controller{
		store:      db.NewLocalDB(),
		publicKey:  pbk,
		privateKey: pvk,
	}

	// TODO: These are not restful at all, but it's the simplest for now. FIX IT!
	r.HandleFunc("/register", c.handleRegister).Methods("POST")
	r.HandleFunc("/registerDevice", c.handleRegisterDevice).Methods("POST")
	r.HandleFunc("/userDevices", c.handleGetUserDevices).Methods("POST")
	r.HandleFunc("/setClipboard", c.handleSetClipboard).Methods("POST")
	r.HandleFunc("/clipboard", c.handleGetClipboard).Methods("POST")

	return r
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

func (c *controller) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Email == "" || req.Password == "" {
		http.Error(w, "both email and password are required", http.StatusBadRequest)
		return
	}

	// TODO: This is obviously just for testing, use Bcrypt or similar for prod.
	err = c.store.PutUser(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := RegisterResponse{Message: "user was successfully registered"}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type RegisterDeviceRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	PublicKey []byte `json:"publicKey"`
}

type RegisterDeviceResponse struct {
	DeviceID string `json:"deviceID"`
	Message  string `json:"message"`
}

// TODO: This should handle devices that are already registered and return the
// existing id.
func (c *controller) handleRegisterDevice(w http.ResponseWriter, r *http.Request) {
	var req RegisterDeviceRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := c.store.GetUser(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: This is obviously just for testing, use Bcrypt or similar for prod.
	if user.PasswordHash != req.Password {
		http.Error(w, "password is not correct for the user", http.StatusUnauthorized)
		return
	}

	deviceId, err := c.store.PutDevice(req.PublicKey, req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := RegisterDeviceResponse{DeviceID: deviceId, Message: "device registered successfully"}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type GetUserDevicesRequest struct {
	FingerPrint `json:",inline"`
}

type GetUserDevicesResponse struct {
	Devices []*api.Device `json:"devices"`
}

func (c *controller) handleGetUserDevices(w http.ResponseWriter, r *http.Request) {
	var authReq AuthenticatedPayload
	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = decryptAuthenticatedPayload[*GetUserDevicesRequest](authReq, c.store, c.privateKey)
	// TODO: verify the request fingerprint. Right now we're just trusting that
	// if it decrypts successfully then we can trust it.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := c.store.GetDeviceUser(authReq.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	devices, err := c.store.GetUserDevices(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := GetUserDevicesResponse{Devices: devices}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type SetClipboardRequest struct {
	FingerPrint `json:",inline"`
	Clipboard   *api.Clipboard `json:"clipboard"`
}

func (c *controller) handleSetClipboard(w http.ResponseWriter, r *http.Request) {
	var authReq AuthenticatedPayload
	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("authReq: ", authReq)

	req, err := decryptAuthenticatedPayload[*SetClipboardRequest](authReq, c.store, c.privateKey)
	// TODO: verify the request fingerprint. Right now we're just trusting that
	// if it decrypts successfully then we can trust it.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("req: ", req)

	user, err := c.store.GetDeviceUser(authReq.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.store.PutClipboard(user.ID, req.Clipboard)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type GetClipboardRequest struct {
	FingerPrint `json:",inline"`
}

type GetClipboardResponse struct {
	Ciphertext      []byte `json:"ciphertext"`
	SenderPublicKey []byte `json:"senderPublicKey"`
}

func (c *controller) handleGetClipboard(w http.ResponseWriter, r *http.Request) {
	var authReq AuthenticatedPayload
	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = decryptAuthenticatedPayload[*GetClipboardRequest](authReq, c.store, c.privateKey)
	// TODO: verify the request fingerprint. Right now we're just trusting that
	// if it decrypts successfully then we can trust it.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user, err := c.store.GetDeviceUser(authReq.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	clip, err := c.store.GetClipboard(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := clip.Payloads[authReq.DeviceID]; !ok {
		http.Error(w, "current clipboard was not produced for this device", http.StatusNotFound)
		return
	}

	res := GetClipboardResponse{SenderPublicKey: clip.SenderPublicKey, Ciphertext: clip.Payloads[authReq.DeviceID]}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
