import requests as rq
import sys
import argparse
import json

def http_get() -> str:
    url: str = "https://www.dwd.de/DWD/warnungen/warnapp/json/warnings.json"
    Headers = {
        'User-Agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/116.0',
        'Cache-Control': 'no-cache',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
        'Host': 'www.dwd.de',
        'Connection': 'keep-alive',
        'Accept-Language': 'de,en-US;q=0.7,en;q=0.3'
        }
    response: rq.Response = rq.get(url, headers=Headers)
    
    if response.status_code == 200:
        return response.text
    else:
        # something bad happend
        exit_unknown(str(response.status_code) + " " + str(response.reason) + "\n" + str(response.headers))

def parse_response(content: str) -> dict:
    # convert dwd's jsonp (what the fck) to beautiful json
    content: str = content.replace('warnWetter.loadWarnings(', '')
    content = content[:-2]
    content_json: dict = json.loads(content)
    return content_json

def exit_unknown(s: str) -> None:
    print(f"UNKNOWN - {s}")
    sys.exit(3)

def exit_critical(s: str) -> None:
    print(s)
    sys.exit(2)

def exit_warning(s: str) -> None:
    print(s)
    sys.exit(1)

def exit_ok(s: str = "No warnings found") -> None:
    print(f"OK - {s}")
    sys.exit(0)


def main() -> None:
    parser = argparse.ArgumentParser(
            prog="check_DWDwarnings",
            description="Get the current warnings from DWD (Deutscher Wetterdienst) by a given cell id. See https://www.dwd.de/DE/leistungen/opendata/help/warnungen/cap_warncellids_csv.csv"
        )
    parser.add_argument("-s", "--station-id", type=int, help="station-id/cell-id", required=True)
    args: argparse.Namespace = parser.parse_args()

    all_warnings: str = http_get()
    station_warning: dict = parse_response(all_warnings)
    if station_warning != None:
        if str(args.station_id) in station_warning["warnings"].keys():
            final_string: str = ""
            highest_level: int = 0
            for warning in station_warning["warnings"][str(args.station_id)]:
                final_string += warning["headline"] + "\n" + warning["description"] + "\n"
                highest_level = max(highest_level, warning["level"])
            if highest_level == 4:
                exit_critical(final_string)
            else:
                exit_warning(final_string)
        else:
            exit_ok()
        
            






if __name__ == "__main__":
    main()
