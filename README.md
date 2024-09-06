Three services.

Implementation of recruitment task.

1. A service that returns a list of JSONs at '/generate/{size}' with the specified size and structure below with random values.
e.g.
{ _type: "Position", _id: 65483214, key: null, name: "Oksywska", fullName: "Oksywska, Poland", iata_airport_code: null, type: “location”, country: "Poland", geo_position:
{ latitude: 51.0855422, longitude: 16.9987442 }, location_id: 756423, inEurope: true,
countryCode: "PL", coreCountry: true, distance: null }

2. A service that retrieves data from the first one and converts it to CSV. The first endpoint that always returns the retrieved data in the format 'type, _id, name, type, latitude, longitude'. The second endpoint returns the retrieved data in the given CSV structure, so we submit in the query 'id, latitude, longitude' and we expect it to return such a result '65483214, 51.0855422, 16.9987442'.
The third endpoint expects the definition of simple mathematical operations in the form of a list, e.g., 'latitude*longitude,sqrt(location_id)' and as a result, will return '3.0052538,869.7258188'

3. A service that performs queries on the second one and displays simple reports regarding performance.
The report should contain information such as the use of processor, memory during the time for each of the previous services, and the time of HTTP queries between services 3->2->1.

Report for 1k, 10k, 100k generated JSONs.

Requirements:

Docker >= 24.0.2
Make >= 4.3

Running project:

Just run make in root directory of repository. Requires docker to run without sudo.