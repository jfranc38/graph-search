package graph_search

// Road Types
const (
	Cycleway      = "cycleway"
	Footway       = "footway"
	Highway       = "highway"
	LivingStreet  = "living_street"
	Motorway      = "motorway"
	MotorwayLink  = "motorway_link"
	Path          = "path"
	Pedestrian    = "pedestrian"
	Primary       = "primary"
	PrimaryLink   = "primary_link"
	Residential   = "residential"
	Road          = "road"
	Secondary     = "secondary"
	SecondaryLink = "secondary_link"
	Service       = "service"
	Steps         = "steps"
	Tertiary      = "tertiary"
	TertiaryLink  = "tertiary_link"
	Track         = "track"
	Trunk         = "trunk"
	TrunkLink     = "trunk_link"
	Unclassified  = "unclassified"
)

// Traffic Control
const (
	Crossing       = "crossing"
	Junction       = "junction"
	Roundabout     = "roundabout"
	TrafficCalming = "traffic_calming"
	TrafficSignals = "traffic_signals"
)

// Road Features
const (
	Intersection  = "intersection"
	TurningCircle = "turning_circle"
	TurningLoop   = "turning_loop"
)

// Direction and Access Control
const (
	No            = "no"
	Oneway        = "oneway"
	Opposite      = "opposite"
	OppositeLane  = "opposite_lane"
	OppositeTrack = "opposite_track"
	Yes           = "yes"
)

// Miscellaneous
const (
	Bicycle  = "bicycle"
	Bike     = "bike"
	Drive    = "drive"
	MaxSpeed = "maxspeed"
)

// SurfaceType constants
const (
	Bricks       = "bricks"
	Cement       = "cement"
	Clay         = "clay"
	Cobblestone  = "cobblestone"
	Compacted    = "compacted"
	Dirt         = "dirt"
	Earth        = "earth"
	FineGravel   = "fine_gravel"
	Grass        = "grass"
	GrassPaver   = "grass_paver"
	Gravel       = "gravel"
	Ground       = "ground"
	Metal        = "metal"
	Mud          = "mud"
	PavingStones = "paving_stones"
	Pebblestone  = "pebblestone"
	Rocky        = "rocky"
	Sand         = "sand"
	Sett         = "sett"
	Surface      = "surface"
	Stone        = "stone"
	Tartan       = "tartan"
	Unpaved      = "unpaved"
	Wood         = "wood"
)

const CellLevel = 30

const (
	AvgSpeedCar              = 40
	AvgSpeedMotor            = 30
	SpeedPenaltyDrive        = 10
	SpeedPenaltyBike         = 5
	SpeedTrafficCalmingDrive = 8
	SpeedTrafficCalmingBike  = 5
)

const (
	MinutesInAnHour    = 60
	MetersInAKilometer = 1000
	KilometersPerMile  = 1.60934
)

var SpeedLimitsSurface = map[string]map[string]float64{
	Drive: {
		Bricks:       60,
		Cement:       80,
		Clay:         30,
		Cobblestone:  30,
		Compacted:    80,
		Dirt:         40,
		Earth:        20,
		FineGravel:   80,
		Grass:        40,
		GrassPaver:   40,
		Gravel:       40,
		Ground:       40,
		Metal:        60,
		Mud:          10,
		PavingStones: 60,
		Pebblestone:  40,
		Rocky:        20,
		Sand:         20,
		Sett:         40,
		Stone:        20,
		Tartan:       40,
		Unpaved:      40,
		Wood:         40,
	},
	Bike: {
		Bricks:       20,
		Cement:       30,
		Clay:         10,
		Cobblestone:  10,
		Compacted:    30,
		Dirt:         15,
		Earth:        5,
		FineGravel:   30,
		Grass:        15,
		GrassPaver:   15,
		Gravel:       15,
		Ground:       15,
		Metal:        20,
		Mud:          5,
		PavingStones: 20,
		Pebblestone:  15,
		Rocky:        5,
		Sand:         5,
		Sett:         15,
		Stone:        5,
		Tartan:       15,
		Unpaved:      15,
		Wood:         15,
	},
}
var SpeedLimitsRoadType = map[string]map[string]float64{
	Drive: {
		LivingStreet:  10,
		Motorway:      89,
		MotorwayLink:  45,
		Primary:       30,
		PrimaryLink:   30,
		Residential:   25,
		Secondary:     49,
		SecondaryLink: 25,
		Service:       15,
		Tertiary:      40,
		TertiaryLink:  20,
		Trunk:         73,
		TrunkLink:     40,
		Unclassified:  25,
	},
	Bike: {
		LivingStreet:  10,
		Motorway:      60,
		MotorwayLink:  30,
		Primary:       30,
		PrimaryLink:   30,
		Residential:   20,
		Secondary:     40,
		SecondaryLink: 20,
		Service:       15,
		Tertiary:      30,
		TertiaryLink:  20,
		Trunk:         50,
		TrunkLink:     30,
		Unclassified:  20,
	},
}
