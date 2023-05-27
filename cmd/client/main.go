package main

func main() {
	clientManager := NewClientManagmentSercieImpl()
	modbusService := NewModbusServiceImpl(clientManager)
	viewController := NewMainModelImpl(modbusService, clientManager)

	MainView(viewController)()
}
