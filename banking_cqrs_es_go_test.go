package main

// func Test_OpenAnAccount(t *testing.T) {
// 	// Given
// 	var eventStore []interface{}
// 	bankingService := banking.NewService(eventStore)
//
// 	accountId, err := uuid.NewV4()
// 	if err != nil {
// 		t.Errorf("Fail err: %+v", err)
// 	}
// 	openAccount, _ := accounts.CreateOpenAccountWithUUID(accountId, "Alex Gemmell")
// 	expectedAccountWasOpened := accounts.CreateAccountWasOpened(accountId.String(), "Alex Gemmell")
//
// 	// When an OpenAccount command is sent
// 	bankingService.HandleCommand(openAccount)
//
// 	// Then an AccountWasOpened event is produced
// 	events := bankingService.eventStore.GetAllEvents()
// 	if len(events) != 1 && events[0] != expectedAccountWasOpened {
// 		t.Fail()
// 	}
// }

// // Move this to a banking service?
// func HandleCommand(account *accounts.OpenAccount) {
// 	// make the argument a generic command type
//
// }
