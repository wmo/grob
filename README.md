# Grob

Grep an OSM Protobuf file. 



# Build





# Usage

Use case: you want to eat a pizza close to Checkpoint Charlie in Berlin.

1) Download the `Berling.osm.pbf` file from [download.geofabrik.de/europe/germany.html](http://download.geofabrik.de/europe/germany.html) (in case the link is no longer valid, start from here: [wiki.openstreetmap.org/wiki/Planet.osm](http://wiki.openstreetmap.org/wiki/Planet.osm)) 

2) Locate checkpoint charlie

    grob berlin-latest.osm.pbf "checkpoint charlie"

    ..
    52.507546, 13.390361,  name:Checkpoint Charlie name:ko:체크포인트 찰리 name:tr:Çarli Kontrol Noktası ..
    ..

3) Check within a radius of 0.5 km of waypoint 52.507546 13.390361 

    grob berlin-latest.osm.pbf -d 0.5 52.507546 13.390361 pizza

    52.505250, 13.393108,  name:Charlotte 1 amenity:restaurant cuisine:italian;pizza #,0.32
    52.506928, 13.392213,  name:Pizza amenity:fast_food #,0.14
    52.506916, 13.395269,  addr:housenumber:9 name:Pepe Pizza amenity:fast_food cuisine:italian .. 



