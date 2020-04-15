package main

type Address struct {
	CEP string `csv:"cep"`
	Latitude string `csv:"latitude"`
	Longitude string `csv:"longitude"`
}
