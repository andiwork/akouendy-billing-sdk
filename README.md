# akouendy-billing-sdk

## Create billing order
### Create Gorm billing transaction table
```
db.AutoMigrate(billing.MigrateModels()...)
```

### Configure and send the request
```
	// create order on billing platform
	billingConfig := billing.Config{Debug: true, Env: "sandbox"}
	billingConfig.AppBaseUrl = viper.GetString("baseurl")
	billing.Init(billingConfig)
	client, err := billing.NewClient(billing.WithRequestBefore(func(ctx context.Context, req *http.Request) error {
		req.Header.Add("X-Billing-Requestid", uuid.New().String())
		req.Header.Add("Authorization", request.HeaderParameter("Authorization"))
		return nil
	}))
	log.Println("===== OrderNewClient ===", err)
	if err == nil {
		request := billing.OrderRequest{}
		request.AppID = viper.GetString("billing-app.id")
		request.PriceID = viper.GetString("billing-app.priceId")
		request.BillingProvider = "orange-money-sn"
		request.CustomerID = userId
		request.CustomerEmail = ""
		request.CustomerFullName = ""
		ctx := context.Background()
		var billingTrx billing.BillingTransaction
		response, billingTrx, err = client.CreateOrder(ctx, settingsId, request)

		// save transactionId and billing ids
		if err == nil {
			repo := billing.BillingRepository{}
			repo.CreateBillingOrder(&billingTrx, utils.GetInstance().GetDB())
		}
	}
```