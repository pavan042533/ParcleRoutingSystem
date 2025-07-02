package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Route struct {
	From     string
	To       string
	Distance int
}

type Parcel struct {
	ID                int
	ShippingPrice     int
	Weight            float64
	TotalDistance     int
	Sender            string
	Receiver          string
	Source            string
	CurrentLoc        string
	Destination       string
	Status            string
	Route             []Route
	CurrentRouteIndex int
}

var parcels []*Parcel

var locations = []string{"Bangalore", "Hyderabad", "Vijayawada", "Chennai", "Delhi"}
var routes = []Route{
	{From: "Bangalore", To: "Hyderabad", Distance: 560},
	{From: "Hyderabad", To: "Vijayawada", Distance: 300},
	{From: "Vijayawada", To: "Chennai", Distance: 520},
	{From: "Hyderabad", To: "Bangalore", Distance: 560},
	{From: "Bangalore", To: "Chennai", Distance: 450},
	{From: "Chennai", To: "Delhi", Distance: 700},
	{From: "Delhi", To: "Bangalore", Distance: 890},
	{From: "Chennai", To: "Bangalore", Distance: 400},
}

func main() {
	fmt.Println("📦 Welcome to Parcel Routing System!")
	mainMenu()
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func CapitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func mainMenu() {
	for {
		fmt.Println("\n=== 📦 Parcel Routing System ===")
		fmt.Println("1. Add Parcel")
		fmt.Println("2. View Parcels")
		fmt.Println("3. Move Parcels")
		fmt.Println("4. Track Parcel by ID") // 👈 NEW OPTION
		fmt.Println("5. Exit")
		fmt.Print("Choice: ")
		choice := readLine()

		switch choice {
		case "1":
			addParcel()
		case "2":
			viewParcels()
		case "3":
			moveParcels()
		case "4":
			trackParcelByID() // 👈 NEW FUNCTION CALL
		case "5":
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("⚠️ Invalid input.")
		}
	}
}

func isValidLocation(name string) bool {
	for _, loc := range locations {
		if loc == name {
			return true
		}
	}
	return false
}

func addParcel() {
	fmt.Println("\nEnter Parcel Details")
	id := len(parcels) + 101
	fmt.Printf("Parcel ID: %v\n", id)

	fmt.Print("Sender: ")
	sender := CapitalizeFirst(readLine())

	fmt.Print("Receiver: ")
	receiver := CapitalizeFirst(readLine())

	fmt.Print("Weight (in KG): ")
	weightInput := readLine()
	weightInt, err := strconv.Atoi(weightInput)
	if err != nil {
		fmt.Println("❌ Invalid weight input.")
		return
	}
	weight := float64(weightInt)

	fmt.Println("Available Locations:", locations)

	fmt.Print("Source Location: ")
	source := CapitalizeFirst(readLine())
	fmt.Print("Destination Location: ")
	destination := CapitalizeFirst(readLine())

	if !isValidLocation(source) || !isValidLocation(destination) {
		fmt.Println("❌ Invalid source or destination.")
		return
	}

	path := findRoute(source, destination)
	if len(path) == 0 {
		fmt.Println("❌ No available route between locations.")
		return
	}

	totalDistance := 0
	for _, r := range path {
		totalDistance += r.Distance
	}

	shippingPrice := calculateShippingPrice(totalDistance, weight)

	if !paymentGateway(shippingPrice, totalDistance, weight, id, sender, receiver, source, destination) {
		fmt.Println("❌ Payment cancelled. Parcel not added.")
		return
	}

	parcel := &Parcel{
		ID:                id,
		ShippingPrice:     shippingPrice,
		Weight:            weight,
		Sender:            sender,
		Receiver:          receiver,
		Source:            source,
		CurrentLoc:        source,
		Destination:       destination,
		Status:            "Pending",
		TotalDistance:     totalDistance,
		Route:             path,
		CurrentRouteIndex: 0,
	}

	parcels = append(parcels, parcel)
	fmt.Println("✅ Parcel added successfully!")
}

func findRoute(source, destination string) []Route {
	var path []Route
	for _, r := range routes {
		if source == r.From && destination == r.To {
			return append(path, r)
		}
	}
	for _, r1 := range routes {
		if r1.From == source {
			for _, r2 := range routes {
				if r1.To == r2.From && r2.To == destination {
					return append(path, r1, r2)
				}
			}
		}
	}
	return []Route{}
}

func calculateShippingPrice(distance int, weight float64) int {
	switch {
	case weight < 5:
		return distance * 2
	case weight < 10:
		return distance * 3
	case weight < 50:
		return distance * 5
	default:
		return distance * 10
	}
}

func viewReceipt(id int, sender, receiver, source, destination string, weight float64, distance, price int) {
	fmt.Println("\n═══════════════════════════════════════════════")
	fmt.Println("               📄 SHIPPING RECEIPT")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Printf("📦 Parcel ID       : %d\n", id)
	fmt.Printf("👤 Sender          : %s\n", sender)
	fmt.Printf("📬 Receiver        : %s\n", receiver)
	fmt.Printf("🏁 From            : %s\n", source)
	fmt.Printf("🏁 To              : %s\n", destination)
	fmt.Printf("⚖️  Weight         : %.2f KG\n", weight)
	fmt.Printf("🛣️  Distance       : %d KM\n", distance)
	fmt.Printf("💸 Shipping Price  : ₹%d\n", price)
	fmt.Println("═══════════════════════════════════════════════")
}

func paymentGateway(price, distance int, weight float64, id int, sender, receiver, source, destination string) bool {
	for {
		fmt.Println("\n=======================================")
		fmt.Printf("💰 Total Shipping Price: ₹%d\n", price)
		fmt.Println("1. View Receipt")
		fmt.Printf("2. Proceed Payment of ₹%d\n", price)
		fmt.Println("3. Cancel Order")
		fmt.Println("=======================================")
		fmt.Print("Enter your choice: ")
		choice := readLine()

		switch choice {
		case "1":
			viewReceipt(id, sender, receiver, source, destination, weight, distance, price)
		case "2":
			fmt.Println("🔐 Processing payment...")
			fmt.Println("✅ Payment Successful.")
			return true
		case "3":
			return false
		default:
			fmt.Println("⚠️ Invalid choice.")
		}
	}
}

func viewParcels() {
	if len(parcels) == 0 {
		fmt.Println("🚫 No parcels to display.")
		return
	}
	fmt.Println("\n══════════════════════════════════════════════════════════════════════════")
	fmt.Println("                         📦 PARCEL STATUS DASHBOARD")
	fmt.Println("══════════════════════════════════════════════════════════════════════════")
	for _, p := range parcels {
		fmt.Printf("ID: %v | From: %v ➡️  %v | Current: %v | Status: %v | Distance: %v KM\n",
			p.ID, p.Source, p.Destination, p.CurrentLoc, p.Status, p.TotalDistance)
	}
	fmt.Printf("🧾 Total Parcels: %d\n", len(parcels))
}

func moveParcels() {
	if len(parcels) == 0 {
		fmt.Println("🚫 No parcels to move.")
		return
	}
	for _, p := range parcels {
		if p.Status == "Delivered" {
			continue
		} else if p.CurrentRouteIndex >= len(p.Route) {
			p.Status = "Delivered"
			continue
		}

		nextHub := p.Route[p.CurrentRouteIndex]
		p.CurrentLoc = nextHub.To
		p.CurrentRouteIndex++

		if p.CurrentLoc == p.Destination {
			p.Status = "Delivered"
		} else if p.CurrentRouteIndex == len(p.Route)-1 {
			p.Status = "Out for Delivery"
		} else {
			p.Status = "In Transit"
		}
	}
	fmt.Println("🚚 Moved all parcels one step forward.")
}

func trackParcelByID() {
	fmt.Print("🔎 Enter Parcel ID to track: ")
	idInput := readLine()
	id, err := strconv.Atoi(idInput)
	if err != nil {
		fmt.Println("❌ Invalid Parcel ID.")
		return
	}

	for _, p := range parcels {
		if p.ID == id {
			fmt.Println("\n📦 Parcel Tracking Details:")
			fmt.Println("═════════════════════════════")
			fmt.Printf("🆔 Parcel ID       : %d\n", p.ID)
			fmt.Printf("👤 Sender          : %s\n", p.Sender)
			fmt.Printf("📬 Receiver        : %s\n", p.Receiver)
			fmt.Printf("🏁 From            : %s\n", p.Source)
			fmt.Printf("📍 Current Location: %s\n", p.CurrentLoc)
			fmt.Printf("🎯 Destination     : %s\n", p.Destination)
			fmt.Printf("🚚 Status          : %s\n", p.Status)
			fmt.Println("═════════════════════════════")
			return
		}
	}
	fmt.Println("❌ Parcel not found.")
}
