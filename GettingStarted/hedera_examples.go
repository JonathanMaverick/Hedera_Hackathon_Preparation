package main

import (
	"fmt"
	"os"

	"github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/joho/godotenv"
)

func main() {

    err := godotenv.Load(".env")
    if err != nil {
        panic(fmt.Errorf("unable to load environment variables from .env file. Error:\n%v", err))
    }

    myAccountId, err := hedera.AccountIDFromString(os.Getenv("MY_ACCOUNT_ID"))
    if err != nil {
        panic(err)
    }

    myPrivateKey, err := hedera.PrivateKeyFromString(os.Getenv("MY_PRIVATE_KEY"))
    if err != nil {
        panic(err)
    }

    fmt.Printf("The account ID is = %v\n", myAccountId)
    fmt.Printf("The private key is = %v\n", myPrivateKey)

	client := hedera.ClientForTestnet()
	client.SetOperator(myAccountId, myPrivateKey)

	client.SetDefaultMaxTransactionFee(hedera.HbarFrom(100, hedera.HbarUnits.Hbar))

	client.SetDefaultMaxQueryPayment(hedera.HbarFrom(50, hedera.HbarUnits.Hbar))

    createAnAccount(myAccountId, *client)
}

func createAnAccount(myAccountId hedera.AccountID,client hedera.Client){
    //2. Generate key for new account
    //Generate new keys for the account you will create
    newAccountPrivateKey, err := hedera.PrivateKeyGenerateEd25519()
    
    if err != nil {
        panic(err)
    }

    newAccountPublicKey := newAccountPrivateKey.PublicKey()

    //3. Create a new account
    //Create new account and assign the public key
    newAccount, err := hedera.NewAccountCreateTransaction().
        SetKey(newAccountPublicKey).
        SetInitialBalance(hedera.HbarFrom(1000, hedera.HbarUnits.Tinybar)).
        Execute(&client)

    if err != nil {
        panic(err)
    }

    //3. Get a new account ID
    //Request the receipt of the transaction
    receipt, err := newAccount.GetReceipt(&client)
    if err != nil {
        panic(err)
    }

    //Get the new account ID from the receipt
    newAccountId := *receipt.AccountID

    fmt.Printf("The new account ID is %v\n", newAccountId)

    //5. Verify the new account balance
    //Create the account balance query
    query := hedera.NewAccountBalanceQuery().SetAccountID(newAccountId)

    //Sign with client operator private key and submit the query to a Hedera network
    accountBalance, err := query.Execute(&client)
    if err != nil {
        panic(err)
    }

    fmt.Println("The account balance for the new account is ", accountBalance.Hbars.AsTinybar())
    transferHbar(myAccountId, newAccountId, client)
}

func transferHbar(myAccountId hedera.AccountID, newAccountId hedera.AccountID, client hedera.Client){
    //Transfer Hbar from your tesnet account to new account
    transaction := hedera.NewTransferTransaction().
        AddHbarTransfer(myAccountId, hedera.HbarFrom(-1000, hedera.HbarUnits.Tinybar)).
        AddHbarTransfer(newAccountId, hedera.HbarFrom(1000, hedera.HbarUnits.Tinybar))

    //Submit the transaction a Hedera network
    txResponse, err := transaction.Execute(&client)
    if err != nil {
        panic(err)
    }

    transferReceipt, err := txResponse.GetReceipt(&client)

    if err != nil {
        panic(err)
    }

    //Get the transaction consensus status
    transactionStatus := transferReceipt.Status
    fmt.Printf("The transaction consesus status is %v\n", transactionStatus)

    queryAccount(newAccountId, client)
}

func queryAccount(newAccountID hedera.AccountID, client hedera.Client){
    //Get the cost of requesting the query
    //Create the query that you want to submit
    balanceQuery := hedera.NewAccountBalanceQuery().SetAccountID(newAccountID)

    //Get the cost of the query
    cost, err := balanceQuery.GetCost(&client)

    if err != nil {
        panic(err)
    }

    fmt.Println("The account balance query cost is ", cost)

    //Get the account balance
    //Check the new account's balance
    newAccountBalanceQuery := hedera.NewAccountBalanceQuery().SetAccountID(newAccountID)

    //Sign with client operator private key and submit the query to a Hedera network
    newAccountBalance, err := newAccountBalanceQuery.Execute(&client)

    if err != nil {
        panic(err)
    }

    fmt.Println("The hbar account balance for this account is", newAccountBalance.Hbars.AsTinybar())
}