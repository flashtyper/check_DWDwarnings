# check_DWDwarnings
Icinga/Nagios plugin which checks the current weather warnings by given cell-id

## Usage: 

```bash
./check_DWDwarnings -s <cell-id>
```

## Example:

```bash
lukas@ubuntu:~$ ./check_DWDwarnings -s 913076001
CRITICAL - Amtliche WARNUNG vor STURMBÖEN
Es treten Sturmböen mit Geschwindigkeiten um 65 km/h (18m/s, 35kn, Bft 8) anfangs aus südwestlicher, später aus westlicher Richtung auf. In Schauernähe sowie in exponierten Lagen muss mit Sturmböen bis 80 km/h (22m/s, 44kn, Bft 9) gerechnet werden.
```


## Installation:
Just move the file into your plugin directory and add a new command and service to your icinga configuration. Don't forget to grant execution permissions. 
Maybe something like this: 

```javascript
object CheckCommand "check_DWDwarnings" {
  command = [ CustomPluginDir + "check_DWDwarnings" ]
  arguments = {
    "-s" = "$cell_ID$"
  }
}
object Service "DWD Warnungen: Stadt Regensburg" {
  check_command = "check_DWDwarnings"
  import "generic-service"
  host_name = "Pi 4"
  vars.cell_ID = "109362000"
}
```

## Cell-IDs
Here you can download a csv-file which contains every cell-id in Germany: https://www.dwd.de/DE/leistungen/opendata/help/warnungen/cap_warncellids_csv.html


