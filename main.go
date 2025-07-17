package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	stripe "github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/account"
	"github.com/stripe/stripe-go/v81/charge"
	"github.com/stripe/stripe-go/v81/checkout/session"
	"github.com/stripe/stripe-go/v81/refund"
)

type request struct {
	Name string `json:"name"`
}

type refundRequest struct {
	Account string `json:"account"`
	Reverse string `json:"reverse"`
	Amount  int64  `json:"amount"`
	AppFee  string `json:"refund_app_fee"`
}

func main() {
	fmt.Println("Starting destination-charge-be")

	r := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"*"}
	r.Use(cors.New(corsConfig))

	r.Handle("POST", "/refund", func(c *gin.Context) {
		fmt.Println("POST /refund")
		var req refundRequest
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			fmt.Printf("Error parsing request: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't parse request",
			})
			return
		}
		reverseRefund, err := strconv.ParseBool(req.Reverse)
		if err != nil {
			fmt.Printf("Error parsing bool from request: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't parse request",
			})
			return
		}
		refundAppFee, err := strconv.ParseBool(req.AppFee)
		if err != nil {
			fmt.Printf("Error parsing bool from request: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't parse request",
			})
			return
		}

		stripe.Key = "sk_test_51C1IhzF1zKN8JIpKx1g8eXCWVlesXbAuBqxmr6iSKiVJoxF7YlkJghvFqSrhz7Jpd4t5oyyL8AgSv8CkP29zqqmp00sPtugADU"
		params := &stripe.RefundParams{
			//TODO FIND CHARGE ID
			Charge:          stripe.String(req.Account),
			ReverseTransfer: stripe.Bool(reverseRefund),
		}
		if req.Amount > 0 {
			params.Amount = stripe.Int64(req.Amount)
		}
		if refundAppFee {
			params.RefundApplicationFee = &refundAppFee
		}

		re, err := refund.New(params)
		if err != nil {
			fmt.Printf("Error creating refund: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't complete request",
			})
			return
		}
		c.JSON(http.StatusOK, re)
	})

	// SERVER ENDPOINT /delete-account
	// Description: endpoint delets the connected account id passed in
	r.Handle("POST", "/delete-account", func(c *gin.Context) {
		fmt.Println("/delete-account called")
		var req request
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			fmt.Printf("Error parsing request: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't parse request",
			})
			return
		}

		stripe.Key = ""

		params := &stripe.AccountParams{}
		fmt.Printf("deleting account: %s\n", req.Name)
		result, err := account.Del(req.Name, params)
		if err != nil {
			fmt.Printf("Error deleting account: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't delete account",
			})
			return
		}
		c.JSON(http.StatusOK, result)
	})

	r.Handle("GET", "/charges", func(c *gin.Context) {
		stripe.Key = "sk_test_51C1IhzF1zKN8JIpKx1g8eXCWVlesXbAuBqxmr6iSKiVJoxF7YlkJghvFqSrhz7Jpd4t5oyyL8AgSv8CkP29zqqmp00sPtugADU"

		params := &stripe.ChargeListParams{}
		// params.Limit = stripe.Int64(7)
		result := charge.List(params)
		c.JSON(http.StatusOK, result)
	})

	// SERVER ENDPOINT /fetch-account
	// Description: endpoint returns the info for the passed in connected account id
	r.Handle("POST", "/fetch-account", func(c *gin.Context) {
		fmt.Println("POST /fetch-account")
		var req request
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			fmt.Printf("Error parsing request: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't parse request",
			})
			return
		}

		fmt.Printf("Fetching account: %s\n", req.Name)
		// stripe.Key = "sk_test_51C1IhzF1zKN8JIpKx1g8eXCWVlesXbAuBqxmr6iSKiVJoxF7YlkJghvFqSrhz7Jpd4t5oyyL8AgSv8CkP29zqqmp00sPtugADU"

		// PRODUCTION ACCOUNT
		stripe.Key = "sk_live_f0zUjoIv8xqPHBcJ9kyMatJP"
		params := &stripe.AccountParams{}
		result, err := account.GetByID(req.Name, params)
		if err != nil {
			fmt.Printf("Error getting account: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't get account",
				"error":   err,
			})
			return

		}
		c.JSON(http.StatusOK, result)
	})

	// SERVER ENDPOINT /update-account
	// Description: endpoint will update the passed in Stripe connected account id
	// currently defaults to Teachable US TP Express Staging Stripe Key
	r.Handle("POST", "/update-account", func(c *gin.Context) {
		fmt.Println("POST /update-account received")

		var req request
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			fmt.Printf("Error parsing request: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't parse request",
			})
			return
		}
		stripe.Key = "sk_test_51C1IhzF1zKN8JIpKx1g8eXCWVlesXbAuBqxmr6iSKiVJoxF7YlkJghvFqSrhz7Jpd4t5oyyL8AgSv8CkP29zqqmp00sPtugADU"
		// UPDATING CONNECTED ACCOUNT DELAYS DAYS
		// params := &stripe.AccountParams{
		// 	Settings: &stripe.AccountSettingsParams{
		// 		Payouts: &stripe.AccountSettingsPayoutsParams{
		// 			Schedule: &stripe.AccountSettingsPayoutsScheduleParams{
		// 				DelayDays: stripe.Int64(5),
		// 			},
		// 		},
		// 	},
		// }

		// UPDATING CONNECTED ACCOUNT TO RECIPIENT SERVICE AGREEMENT
		params := &stripe.AccountParams{
			TOSAcceptance: &stripe.AccountTOSAcceptanceParams{
				ServiceAgreement: stripe.String("recipient"),
				Date:             stripe.Int64(1609798905),
				IP:               stripe.String("0.0.0.0"),
			},
		}

		acct, err := account.Update(req.Name, params)
		if err != nil {
			fmt.Printf("Error getting account: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't get account",
				"error":   err,
			})
			return

		}
		c.JSON(http.StatusOK, acct)
	})

	// SERVER ENDPOINT /create-checkout-session
	// Description: Uses Teachable US TP Express Staging Stripe Key and starts ;a
	// Checkout session returning the ClientSecret
	r.Handle("POST", "/create-checkout-session", func(c *gin.Context) {
		fmt.Println("POST /create-checkout-session received")

		var req request
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			fmt.Printf("Error parsing request: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't parse request",
			})
			return
		}
		stripe.Key = "sk_test_51C1IhzF1zKN8JIpKx1g8eXCWVlesXbAuBqxmr6iSKiVJoxF7YlkJghvFqSrhz7Jpd4t5oyyL8AgSv8CkP29zqqmp00sPtugADU"
		/*
			NZ - acct_1JYGw72V1ZrQbUTv
			CANADA - acct_1IBQYU2X2GsXY6qD
			JAPAN - acct_1QKMqZFRTWGEiw6J
			INDIA - acct_1M9dq3FJRujqPjso
			EU - FRANCE - acct_1IBjzd2Vn2gzlzfr
			US - acct_1FOVnBESiLUvGpSq
		*/

		destination_account := "acct_1IBjzd2Vn2gzlzfr"
		params := &stripe.CheckoutSessionParams{
			LineItems: []*stripe.CheckoutSessionLineItemParams{
				&stripe.CheckoutSessionLineItemParams{
					PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
						Currency: stripe.String(string(stripe.CurrencyEUR)),
						// Currency: stripe.String(string(stripe.CurrencyIDR)),
						// Currency: stripe.String(string(stripe.CurrencyCAD)),
						// Currency: stripe.String(string(stripe.CurrencyEUR)),
						// Currency: stripe.String(string(stripe.CurrencyUSD)),
						ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
							Name: stripe.String("T-shirt"),
						},
						UnitAmount: stripe.Int64(1200),
					},
					Quantity: stripe.Int64(1),
				},
			},
			PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
				// OnBehalfOf:           stripe.String(destination_account),
				ApplicationFeeAmount: stripe.Int64(111),
				TransferData: &stripe.CheckoutSessionPaymentIntentDataTransferDataParams{
					Destination: stripe.String(destination_account),
				},
			},
			Mode:      stripe.String(string(stripe.CheckoutSessionModePayment)),
			UIMode:    stripe.String(string(stripe.CheckoutSessionUIModeEmbedded)),
			ReturnURL: stripe.String("http://localhost:3000.com/checkout/return?session_id={CHECKOUT_SESSION_ID}"),
		}

		result, err := session.New(params)
		if err != nil {
			fmt.Printf("Error creating session: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't create session",
			})
			return
		}

		jsn, err := json.Marshal(result)
		if err != nil {
			fmt.Printf("Error  marshalling json: %+v\n", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "couldn't create session",
			})
			return
		}
		fmt.Printf("result: %+v", result)
		fmt.Println(string(jsn))

		c.JSON(http.StatusOK, result.ClientSecret)
	})

	r.GET("/ping", func(c *gin.Context) {
		// fmt.Println("Starting destination-charge-be")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
