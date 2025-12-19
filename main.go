package main

func main() {
    as111 := map[string]map[string]string {
        "R1": {
            "GigabitEthernet 1/0": "2001:1:100::1/64",
        },
        "R2": {
            "GigabitEthernet 1/0": "2001:1:100::2/64",
            "GigabitEthernet 2/0": "2001:1:101::1/64",
        },
    }

	// Create a .cfg file for router named R1 with data dict
    writeConfig("R1", as111, "RIP")
}
