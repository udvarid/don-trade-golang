package user

import (
	"math/rand"
	"slices"
	"time"
)

var VORNAMES = []string{"Avaricious", "Covetous", "Rapacious", "Grasping", "Voracious", "Acquisitive", "Insatiable",
	"Selfish", "Materialistic", "Gluttonous", "Miserly", "Possessive", "Hoarding", "Mercenary", "Avarice", "Opportunistic",
	"Egotistic", "Envious", "Parasitic", "Predatory", "Consuming", "Glutton", "Stingy", "Excessive", "Monopolizing",
	"Self-Centered", "Unscrupulous", "Greedy-Guts", "Exploitive", "Hoarding", "Cupidity",
	"Piggy", "Voracious", "Avaricious", "Opportunist", "Predatory",
	"Voracious", "Grabbing", "Gluttonous", "Greedy-Pig", "Inordinate", "Insatiate", "Covetous", "Rapacious", "Fortunate",
	"Blessed", "Favored", "Auspicious", "Fortunate", "Prosperous", "Serendipitous", "Charmed", "Advantageous", "Providential",
	"Opportune", "Fortuitous", "Favorable", "Timely", "Auspicious", "Prosperous", "Charmed", "Successful", "Serendipitous",
	"Advantageous", "Felicitous", "Fortunate", "Providential", "Propitious", "Smiling", "Promising", "Bright", "Prosperous",
	"Rosy", "Golden", "Hopeful", "Encouraging", "Heartening", "Windfall", "Boon", "Break", "Shot", "Happening", "Destined",
	"Unexpected", "Golden", "Auspicious", "Diligent", "Industrious", "Conscientious", "Persistent", "Assiduous",
	"Dedicated", "Tenacious", "Committed", "Laborious", "Earnest", "Persevering", "Tireless", "Vigorous", "Relentless",
	"Studious", "Zealous", "Active", "Busy", "Driven", "Focused", "Serious", "Steady", "Unflagging", "Untiring", "Unwearying",
	"Unrelenting", "Meticulous", "Resolute", "Methodical", "Attentive", "Purposeful", "Engaged", "Hardworking", "Punctual",
	"Scrupulous", "Spirited", "Devoted", "Exacting", "Energetic", "Productive", "Striving", "Determined", "Ambitious",
	"Intelligent", "Smart", "Shrewd", "Ingenious", "Resourceful", "Quick-Witted", "Cunning", "Savvy", "Astute", "Crafty",
	"Bright", "Brilliant", "Sagacious", "Perceptive", "Witty", "Keen", "Insightful", "Sharp", "Prudent", "Tactful",
	"Skillful", "Strategic", "Adroit", "Nimble", "Slick", "Deft", "Artful", "Calculating", "Discerning", "Innovative",
	"Inventive", "Original", "Wise", "Proficient", "Competent", "Brainy", "Capable", "Genius", "Brainy", "Foxy", "Quick",
	"Sharp", "Intuitive", "Strategic", "Artful", "Genius", "Adept"}
var FAMILYNAMES = []string{"Merchant", "Dealer", "Vendor", "Broker", "Retailer", "Wholesaler", "Supplier",
	"Distributor", "Businessperson", "Entrepreneur", "Shopkeeper", "Seller", "Peddler", "Jobber", "Middleman",
	"Negotiator", "Marketer", "Hawker", "Auctioneer", "Agent", "Representative", "Stockbroker", "Financier",
	"Investor", "Salesperson", "Speculator", "Commercialist", "Negociant", "Promoter", "Tradesman", "Tycoon",
	"Concessionaire", "Businessman", "Economist", "Realtor", "Shopman", "Mercantilist", "Buyer", "Procurement",
	"Venturer", "Exchanger", "Purveyor", "Monger", "Salesman", "Vender", "Provider", "Distributer", "Marketeer"}

func GetRandomUniqueName(names []string) string {
	unique := false
	var result string
	tries := 0
	for {
		result = VORNAMES[getRandomNumber(len(VORNAMES))] + " " + FAMILYNAMES[getRandomNumber(len(FAMILYNAMES))]
		unique = !slices.Contains(names, result)
		tries++
		if unique || tries > 1000 {
			break
		}

	}
	return result
}

func getRandomNumber(max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rng := rand.New(r)
	return rng.Intn(max)
}
