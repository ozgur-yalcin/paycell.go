package paycell

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	Application = "PAYCELLTEST"
	Password    = "PaycellTestPassword"
	StoreKey    = "PAYCELL12345"
	Merchant    = "9998"
	EulaID      = "17"
	Prefix      = "666"
	Endpoint    = map[string]string{
		"PROD":       "https://tpay.turkcell.com.tr/tpay/provision/services/restful/getCardToken",
		"TEST":       "https://tpay-test.turkcell.com.tr/tpay/provision/services/restful/getCardToken",
		"PROD_TOKEN": "https://epayment.turkcell.com.tr/paymentmanagement/rest/getCardTokenSecure",
		"PROD_FORM":  "https://epayment.turkcell.com.tr/paymentmanagement/rest/threeDSecure",
		"TEST_TOKEN": "https://omccstb.turkcell.com.tr/paymentmanagement/rest/getCardTokenSecure",
		"TEST_FORM":  "https://omccstb.turkcell.com.tr/paymentmanagement/rest/threeDSecure",
	}
)

type any = interface{}

type API struct {
	Mode     string
	MSisdn   string
	ClientIP string
	Amount   string
	Currency string
}

type Request struct {
	CardToken struct {
		Header     RequestHeader `json:"requestHeader,omitempty"`
		CardNumber any           `json:"creditCardNo,omitempty"`
		CardMonth  any           `json:"expireDateMonth,omitempty"`
		CardYear   any           `json:"expireDateYear,omitempty"`
		CardCode   any           `json:"cvcNo,omitempty"`
		HashData   any           `json:"hashData,omitempty"`
	}
	Provision struct {
		Header        RequestHeader `json:"requestHeader,omitempty"`
		MSisdn        any           `json:"msisdn,omitempty"`
		MerchantCode  any           `json:"merchantCode,omitempty"`
		CardId        any           `json:"cardId,omitempty"`
		CardToken     any           `json:"cardToken,omitempty"`
		RefNo         any           `json:"referenceNumber,omitempty"`
		OriginalRefNo any           `json:"originalReferenceNumber,omitempty"`
		Amount        any           `json:"amount,omitempty"`
		PointAmount   any           `json:"pointAmount,omitempty"`
		Currency      any           `json:"currency,omitempty"`
		Installment   any           `json:"installmentCount,omitempty"`
		PaymentType   any           `json:"paymentType,omitempty"`
		AcquirerBank  any           `json:"acquirerBankCode,omitempty"`
		ThreeDSession any           `json:"threeDSessionId,omitempty"`
		Pin           any           `json:"pin,omitempty"`
	}
	Refund struct {
		Header        RequestHeader `json:"requestHeader,omitempty"`
		MSisdn        any           `json:"msisdn,omitempty"`
		MerchantCode  any           `json:"merchantCode,omitempty"`
		Amount        any           `json:"amount,omitempty"`
		Currency      any           `json:"currency,omitempty"`
		RefNo         any           `json:"referenceNumber,omitempty"`
		OriginalRefNo any           `json:"originalReferenceNumber,omitempty"`
	}
	Reverse struct {
		Header        RequestHeader `json:"requestHeader,omitempty"`
		MSisdn        any           `json:"msisdn,omitempty"`
		MerchantCode  any           `json:"merchantCode,omitempty"`
		RefNo         any           `json:"referenceNumber,omitempty"`
		OriginalRefNo any           `json:"originalReferenceNumber,omitempty"`
	}
	ThreeDSession struct {
		Header       RequestHeader `json:"requestHeader,omitempty"`
		MSisdn       any           `json:"msisdn,omitempty"`
		MerchantCode any           `json:"merchantCode,omitempty"`
		CardId       any           `json:"cardId,omitempty"`
		CardToken    any           `json:"cardToken,omitempty"`
		Installment  any           `json:"installmentCount,omitempty"`
		Amount       any           `json:"amount,omitempty"`
		Currency     any           `json:"currency,omitempty"`
		RefNo        any           `json:"referenceNumber,omitempty"`
		Target       any           `json:"target,omitempty"`
		Transaction  any           `json:"transactionType,omitempty"`
	}
	ThreeDResult struct {
		Header        RequestHeader `json:"requestHeader,omitempty"`
		MSisdn        any           `json:"msisdn,omitempty"`
		MerchantCode  any           `json:"merchantCode,omitempty"`
		RefNo         any           `json:"referenceNumber,omitempty"`
		ThreeDSession any           `json:"threeDSessionId,omitempty"`
	}
	ThreeDForm struct {
		ThreeDSession  any `form:"threeDSessionId,omitempty"`
		CallbackUrl    any `form:"callbackurl,omitempty"`
		IsPoint        any `form:"isPoint,omitempty"`
		IsPost3DResult any `form:"isPost3DResult,omitempty"`
	}
	PaymentMethods struct {
		Header RequestHeader `json:"requestHeader,omitempty"`
		MSisdn any           `json:"msisdn,omitempty"`
	}
	MobilePayment struct {
		Header RequestHeader `json:"requestHeader,omitempty"`
		MSisdn any           `json:"msisdn,omitempty"`
		EulaID any           `json:"eulaID,omitempty"`
	}
	OTP struct {
		Header   RequestHeader `json:"requestHeader,omitempty"`
		MSisdn   any           `json:"msisdn,omitempty"`
		Amount   any           `json:"amount,omitempty"`
		Currency any           `json:"currency,omitempty"`
		RefNo    any           `json:"referenceNumber,omitempty"`
		OTP      any           `json:"otp,omitempty"`
		Token    any           `json:"token,omitempty"`
	}
}

type Response struct {
	CardToken struct {
		Header    ResponseHeader `json:"responseHeader,omitempty"`
		CardToken any            `json:"cardToken,omitempty"`
		HashData  any            `json:"hashData,omitempty"`
	}
	Provision struct {
		Header       ResponseHeader `json:"responseHeader,omitempty"`
		OrderId      any            `json:"orderId,omitempty"`
		OrderDate    any            `json:"reconciliationDate,omitempty"`
		ApprovalCode any            `json:"approvalCodeo,omitempty"`
		AcquirerBank any            `json:"acquirerBankCode,omitempty"`
		IssuerBank   any            `json:"issuerBankCode,omitempty"`
	}
	ThreeDSession struct {
		Header        ResponseHeader `json:"responseHeader,omitempty"`
		ThreeDSession any            `json:"threeDSessionId,omitempty"`
	}
	ThreeDResult struct {
		CurrentStep    any `json:"currentStep,omitempty"`
		MdErrorMessage any `json:"mdErrorMessage,omitempty"`
		MdStatus       any `json:"mdStatus,omitempty"`
		Operation      struct {
			Result      any `json:"threeDResult,omitempty"`
			Description any `json:"threeDResultDescription,omitempty"`
		} `json:"threeDOperationResult,omitempty"`
	}
	PaymentMethods struct {
		Header   ResponseHeader `json:"responseHeader,omitempty"`
		EulaID   any            `json:"eulaID,omitempty"`
		CardList []struct {
			CardBrand         any  `json:"cardBrand,omitempty"`
			CardId            any  `json:"cardId,omitempty"`
			CardType          any  `json:"cardType,omitempty"`
			MaskedCardNo      any  `json:"maskedCardNo,omitempty"`
			Alias             any  `json:"alias,omitempty"`
			ActivationDate    any  `json:"activationDate,omitempty"`
			IsDefault         bool `json:"isDefault,omitempty"`
			IsExpired         bool `json:"isExpired,omitempty"`
			ShowEulaId        bool `json:"showEulaId,omitempty"`
			IsThreeDValidated bool `json:"isThreeDValidated,omitempty"`
			IsOTPValidated    bool `json:"isOTPValidated,omitempty"`
		} `json:"cardList,omitempty"`
		MobilePayment *struct {
			EulaId         any  `json:"eulaId,omitempty"`
			EulaUrl        any  `json:"eulaUrl,omitempty"`
			SignedEulaId   any  `json:"signedEulaId,omitempty"`
			StatementDate  any  `json:"statementDate,omitempty"`
			Limit          any  `json:"limit,omitempty"`
			MaxLimit       any  `json:"maxLimit,omitempty"`
			RemainingLimit any  `json:"remainingLimit,omitempty"`
			IsDcbOpen      bool `json:"isDcbOpen,omitempty"`
			IsEulaExpired  bool `json:"isEulaExpired,omitempty"`
		} `json:"mobilePayment,omitempty"`
	}
	MobilePayment struct {
		Header ResponseHeader `json:"responseHeader,omitempty"`
	}
	OTP struct {
		Header     ResponseHeader `json:"responseHeader,omitempty"`
		Token      any            `json:"token,omitempty"`
		ExpireDate any            `json:"expireDate,omitempty"`
		RetryCount any            `json:"remainingRetryCount,omitempty"`
	}
}

type RequestHeader struct {
	ApplicationName     string `json:"applicationName,omitempty"`
	ApplicationPwd      string `json:"applicationPwd,omitempty"`
	ClientIPAddress     string `json:"clientIPAddress,omitempty"`
	TransactionDateTime string `json:"transactionDateTime,omitempty"`
	TransactionId       string `json:"transactionId,omitempty"`
}

type ResponseHeader struct {
	ResponseCode        string `json:"responseCode,omitempty"`
	ResponseDescription string `json:"responseDescription,omitempty"`
	ResponseDateTime    string `json:"responseDateTime,omitempty"`
	TransactionId       string `json:"transactionId,omitempty"`
}

func SHA256(data string) (hash string) {
	h := sha256.New()
	h.Write([]byte(data))
	hash = base64.StdEncoding.EncodeToString(h.Sum(nil))
	return hash
}

func Random(n int) string {
	const alphanum = "0123456789"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func Api(msisdn string) (*API, *Request) {
	api := new(API)
	api.MSisdn = msisdn
	request := new(Request)
	return api, request
}

func (api *API) SetMode(mode string) {
	api.Mode = mode
}

func (api *API) SetIPAddress(ip string) {
	api.ClientIP = ip
}

func (api *API) SetAmount(total string, currency string) {
	api.Amount = strings.ReplaceAll(total, ".", "")
	api.Currency = currency
}

func (request *Request) SetCardNumber(number string) {
	request.CardToken.CardNumber = number
}

func (request *Request) SetCardExpiry(month, year string) {
	request.CardToken.CardMonth = month
	request.CardToken.CardYear = year
}

func (request *Request) SetCardCode(code string) {
	request.CardToken.CardCode = code
}

func (api *API) HashResponse(header ResponseHeader, cardToken string) string {
	hashdata := SHA256(strings.ToUpper(Application + header.TransactionId + header.ResponseDateTime + header.ResponseCode + cardToken + StoreKey + SHA256(strings.ToUpper(Password+Application))))
	return hashdata
}

func (api *API) Auth() (response Response) {
	apiurl := Endpoint[api.Mode] + "/provision/"
	request := new(Request)
	request.Provision.Header.ClientIPAddress = api.ClientIP
	request.Provision.Header.ApplicationName = Application
	request.Provision.Header.ApplicationPwd = Password
	request.Provision.Header.TransactionDateTime = strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")
	request.Provision.Header.TransactionId = Random(20)
	request.Provision.MSisdn = api.MSisdn
	request.Provision.MerchantCode = Merchant
	request.Provision.RefNo = Prefix + fmt.Sprintf("%v", request.Provision.Header.TransactionDateTime)
	request.Provision.Amount = api.Amount
	request.Provision.Currency = api.Currency
	request.Provision.PaymentType = "SALE"
	postdata, _ := json.Marshal(request.Provision)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return response
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	decoder.Decode(&response.Provision)
	return response
}

func (api *API) PreAuth() (response Response) {
	apiurl := Endpoint[api.Mode] + "/provision/"
	request := new(Request)
	request.Provision.Header.ClientIPAddress = api.ClientIP
	request.Provision.Header.ApplicationName = Application
	request.Provision.Header.ApplicationPwd = Password
	request.Provision.Header.TransactionDateTime = strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")
	request.Provision.Header.TransactionId = Random(20)
	request.Provision.MSisdn = api.MSisdn
	request.Provision.MerchantCode = Merchant
	request.Provision.RefNo = Prefix + fmt.Sprintf("%v", request.Provision.Header.TransactionDateTime)
	request.Provision.Amount = api.Amount
	request.Provision.Currency = api.Currency
	request.Provision.PaymentType = "PREAUTH"
	postdata, _ := json.Marshal(request.Provision)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return response
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	decoder.Decode(&response.Provision)
	return response
}

func (api *API) PostAuth() (response Response) {
	apiurl := Endpoint[api.Mode] + "/provision/"
	request := new(Request)
	request.Provision.Header.ClientIPAddress = api.ClientIP
	request.Provision.Header.ApplicationName = Application
	request.Provision.Header.ApplicationPwd = Password
	request.Provision.Header.TransactionDateTime = strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")
	request.Provision.Header.TransactionId = Random(20)
	request.Provision.MSisdn = api.MSisdn
	request.Provision.MerchantCode = Merchant
	request.Provision.RefNo = Prefix + fmt.Sprintf("%v", request.Provision.Header.TransactionDateTime)
	request.Provision.Amount = api.Amount
	request.Provision.Currency = api.Currency
	request.Provision.PaymentType = "POSTAUTH"
	postdata, _ := json.Marshal(request.Provision)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return response
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	decoder.Decode(&response.Provision)
	return response
}

func (api *API) ThreeDSession() (response Response) {
	apiurl := Endpoint[api.Mode] + "/getThreeDSession/"
	request := new(Request)
	request.ThreeDSession.Header.ClientIPAddress = api.ClientIP
	request.ThreeDSession.Header.ApplicationName = Application
	request.ThreeDSession.Header.ApplicationPwd = Password
	request.ThreeDSession.Header.TransactionDateTime = strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")
	request.ThreeDSession.Header.TransactionId = Random(20)
	request.ThreeDSession.MSisdn = api.MSisdn
	request.ThreeDSession.MerchantCode = Merchant
	request.ThreeDSession.RefNo = Prefix + fmt.Sprintf("%v", request.ThreeDSession.Header.TransactionDateTime)
	request.ThreeDSession.Amount = api.Amount
	request.ThreeDSession.Currency = api.Currency
	postdata, _ := json.Marshal(request.ThreeDSession)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return response
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	decoder.Decode(&response.ThreeDSession)
	return response
}

func (api *API) ThreeDResult(session interface{}) (response Response) {
	apiurl := Endpoint[api.Mode] + "/getThreeDSessionResult/"
	request := new(Request)
	request.ThreeDResult.Header.ClientIPAddress = api.ClientIP
	request.ThreeDResult.Header.ApplicationName = Application
	request.ThreeDResult.Header.ApplicationPwd = Password
	request.ThreeDResult.Header.TransactionDateTime = strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")
	request.ThreeDResult.Header.TransactionId = Random(20)
	request.ThreeDResult.MSisdn = api.MSisdn
	request.ThreeDResult.MerchantCode = Merchant
	request.ThreeDResult.RefNo = Prefix + fmt.Sprintf("%v", request.ThreeDResult.Header.TransactionDateTime)
	if session != nil {
		request.ThreeDResult.ThreeDSession = session
	}
	postdata, _ := json.Marshal(request.ThreeDResult)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return response
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	decoder.Decode(&response.ThreeDResult)
	return response
}

func (api *API) CardToken(request *Request) (response Response) {
	apiurl := Endpoint[api.Mode+"_TOKEN"]
	request.CardToken.HashData = SHA256(strings.ToUpper(Application + request.CardToken.Header.TransactionId + request.CardToken.Header.TransactionDateTime + StoreKey + SHA256(strings.ToUpper(Password+Application))))
	postdata, _ := json.Marshal(request.CardToken)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return response
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	decoder.Decode(&response.CardToken)
	return response
}

func (api *API) GetPaymentMethods() (response Response) {
	apiurl := Endpoint[api.Mode] + "/getPaymentMethods/"
	request := new(Request)
	request.PaymentMethods.MSisdn = api.MSisdn
	request.PaymentMethods.Header.ClientIPAddress = api.ClientIP
	request.PaymentMethods.Header.ApplicationName = Application
	request.PaymentMethods.Header.ApplicationPwd = Password
	request.PaymentMethods.Header.TransactionDateTime = strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")
	request.PaymentMethods.Header.TransactionId = Random(20)
	postdata, _ := json.Marshal(request.PaymentMethods)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return response
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	decoder.Decode(&response.PaymentMethods)
	return response
}

func (api *API) OpenMobilePayment(eula interface{}) (response Response) {
	apiurl := Endpoint[api.Mode] + "/openMobilePayment/"
	request := new(Request)
	request.MobilePayment.Header.ClientIPAddress = api.ClientIP
	request.MobilePayment.Header.ApplicationName = Application
	request.MobilePayment.Header.ApplicationPwd = Password
	request.MobilePayment.Header.TransactionDateTime = strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")
	request.MobilePayment.Header.TransactionId = Random(20)
	request.MobilePayment.MSisdn = api.MSisdn
	if eula != nil {
		request.MobilePayment.EulaID = eula
	}
	postdata, _ := json.Marshal(request.MobilePayment)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return response
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	decoder.Decode(&response.MobilePayment)
	return response
}

func (api *API) SendOTP() (response Response) {
	apiurl := Endpoint[api.Mode] + "/sendOTP/"
	request := new(Request)
	request.OTP.Header.ClientIPAddress = api.ClientIP
	request.OTP.Header.ApplicationName = Application
	request.OTP.Header.ApplicationPwd = Password
	request.OTP.Header.TransactionDateTime = strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")
	request.OTP.Header.TransactionId = Random(20)
	request.OTP.MSisdn = api.MSisdn
	request.OTP.RefNo = Random(20)
	request.OTP.Amount = api.Amount
	postdata, _ := json.Marshal(request.OTP)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return response
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	decoder.Decode(&response.OTP)
	return response
}

func (api *API) ValidateOTP(token, otp interface{}) (response Response) {
	apiurl := Endpoint[api.Mode] + "/validateOTP/"
	request := new(Request)
	request.OTP.Header.ClientIPAddress = api.ClientIP
	request.OTP.Header.ApplicationName = Application
	request.OTP.Header.ApplicationPwd = Password
	request.OTP.Header.TransactionDateTime = strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "")
	request.OTP.Header.TransactionId = Random(20)
	request.OTP.MSisdn = api.MSisdn
	request.OTP.RefNo = Random(20)
	request.OTP.Amount = api.Amount
	if token != nil {
		request.OTP.Token = token
	}
	if otp != nil {
		request.OTP.OTP = otp
	}
	postdata, _ := json.Marshal(request.OTP)
	cli := new(http.Client)
	req, err := http.NewRequest("POST", apiurl, bytes.NewReader(postdata))
	if err != nil {
		log.Println(err)
		return response
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		log.Println(err)
		return response
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	decoder.UseNumber()
	decoder.Decode(&response.OTP)
	return response
}
