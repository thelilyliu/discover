var dirDisp

function initMap () {
    map = new google.maps.Map(document.getElementById('map'), {
        center: { lat: -34.397, lng: 150.644 },
        zoom: 16
    });

    infoWindow = new google.maps.InfoWindow;
    
    dirDisp = new google.maps.DirectionsRenderer({
        map: map
    })

    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(function (position) {
            pos = {
                lat: position.coords.latitude,
                lng: position.coords.longitude
            };

        //    infoWindow.setPosition(pos);
        //    infoWindow.setContent('Location found.');
        //    infoWindow.open(map);
            map.setCenter(pos);
            
            //plotPoints(point_set);

        }, function () {
            handleLocationError(true, infoWindow, map.getCenter());
        });
    } else {
        // Browser doesn't support Geolocation
        handleLocationError(false, infoWindow, map.getCenter());
    }
}

function handleLocationError(browserHasGeolocation, infoWindow, pos) {
    infoWindow.setPosition(pos);
    infoWindow.setContent(browserHasGeolocation ?
        'Error: The Geolocation service failed.' :
        'Error: Your browser doesn\'t support geolocation.');
    infoWindow.open(map);
}

var one = [
    {
        "location": {
            "lat": 43.661298141724274,
            "lng": -79.4013948957703
        },
        "stopover": true
    },
    {
        "location": {
            "lat": 43.6579468,
            "lng": -79.4001475
        },
        "stopover": true
    }
]

function plotPoints (waypointsArr) {
    var pos = {lat: 43.6595063, lng: -79.397758}
    var waypointsArrFinal = []
    var final_point;
    for (var i = 0; i < waypointsArr.length; i++) {
        var waypoint = {
            "location" : {
                "lat" : waypointsArr[i].latitude,
                "lng" : waypointsArr[i].longitude
            },
            "stopover": true
        }
        waypointsArrFinal.push(waypoint)
    }
    var request = {
        // from: Blackpool to: Preston to: Blackburn
        origin: pos,
        destination: pos,
        waypoints: waypointsArrFinal,
        optimizeWaypoints: true,
        travelMode: google.maps.DirectionsTravelMode.WALKING
    };


    var service = new google.maps.DirectionsService()

    service.route(request, function (response, status) {
        if (status == google.maps.DirectionsStatus.OK) {
          //  directionsDisplay.setDirection(response)
            var route = response.routes[0]
        }
        dirDisp.setDirections(response)
    })

    setInterval(function() {
        
    },3000)
}


function computeTotalDistance(result) {
    var total_dist = 0;
    var total_time = 0;
    var myroute = result.routes[0];
    
    /*for (var x = 0; x < result.request.waypoints.length; x++) {
        console.log(result.request.waypoints[x].location.location.lat(), result.request.waypoints[x].location.location.lng())
    }*/

    for (var i = 0; i < myroute.legs.length; i++) {
        
        total_dist += myroute.legs[i].distance.value;
        total_time += myroute.legs[i].duration.value;
    }
    var data = [total_time, total_dist]
    console.log(data)
    return data
}

