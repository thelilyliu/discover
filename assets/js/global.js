var map
var points
var data
var index = 0
var userID = "5a63f03d56cd53a171000001"

$(document).ready(function () {
    $('select').material_select()

    eventHandler()
})

function eventHandler() {
    $('#option').on('click', 'a.send', function () {
        var option = $('form select option:selected').val()

        sendOption(option)

        $(this).addClass('disabled')
    })

    $('#upload-file').change(function () {
        var input = document.querySelector('#upload-file')

        uploadImage(input)
    })
}

/*
  ========================================
  Initialization
  ========================================
*/

function displayMap() {
    $('#map').css('display', 'block')
    initMap()
    displayNext(points[index], 0, 0, 0)
}

function displayNext(point, score, dist, time) {
    var $this = $('#next')
    $this.find('.data-address').text(point.address)

    var $chipWrapper = $this.find('.chip-wrapper')
    $chipWrapper.empty()

    for (var i = 0; i < point.keywords.length; i++) {
        var chip = '<div class="chip">' + point.keywords[i] + '</div>'
        $chipWrapper.append(chip)
    }

    $this.find('.data-score').text(' ' + score + ' pts')
    $this.find('.data-dist').text(' ' + dist + ' km')
    $this.find('.data-time').text(' ' + time + ' min')

    $('#next').css('display', 'block')
    $('#upload').css('display', 'block')
}

/*
  ========================================
  Ajax
  ========================================
*/

function sendOption(option) {
    $.ajax({
        type: 'GET',
        url: '/sendOption/' + userID + '/' + option,
        dataType: 'json',
        cache: false
    }).done(function (json, textStatus, jqXHr) {
        console.log('done')

        points = json.points
        displayMap()

        var pos = { lat: 43.6595063, lng: -79.397758 }
        var waypointsArrFinal = []

        for (var i = 0; i < points.length; i++) {
            console.log(i)
            var waypoint = {
                "location": {
                    "lat": points[i].latitude,
                    "lng": points[i].longitude
                },
                "stopover": true
            }
            waypointsArrFinal.push(waypoint)
        }
        
        var request = {
            // from: Blackpool to: Preston to: Blackburn
            origin: {"lat": points[currentPosition].latitude,
                "lng": points[currentPosition].longitude},
            destination: pos,
            waypoints: waypointsArrFinal,
            optimizeWaypoints: true,
            travelMode: google.maps.DirectionsTravelMode.WALKING
        }

        console.log(request)
        var service = new google.maps.DirectionsService()
        service.route(request, function (response, status) {
            if (status == google.maps.DirectionsStatus.OK) {
                //  directionsDisplay.setDirection(response)
                var route = response.routes[0]
                var data = computeTotalDistance(response)
            }
        })

        plotPoints(points)
    }).fail(function (jqXHr, textStatus, errorThrown) {
        console.log('fail')
    }).always(function () { })
}

function uploadImage(input) {
    var xhr = new XMLHttpRequest()
    var url = '/postImage'
    var fd = new FormData()

    fd.append('uploadFile', input.files[0])
    xhr.open('POST', url, true)

    xhr.onreadystatechange = function () {
        if (xhr.readyState == 4 && xhr.status == 200) {
            if (xhr.responseText != 'fail') {
                pointCheckResult()
            } else {
                console.log('fail')
            }
        }
    }

    xhr.send(fd)
}

var currentPosition = 0;

function pointCheckResult() {
    $.ajax({
        type: 'POST',
        url: '/pointCheckResult/' + points[index].pointID,
        data: JSON.stringify(points[index]),
        contentType: 'application/json; charset=utf-8',
        dataType: 'json',
        cache: false
    }).done(function (json, textStatus, jqXHr) {
        if (json) {
            points2 = points.slice(currentPosition, points.length-1)
            console.log(points2, currentPosition)

            var pos = { lat: 43.6595063, lng: -79.397758 }
            var waypointsArrFinal = []

            for (var i = currentPosition; i < points2.length; i++) {
                console.log(i)
                var waypoint = {
                    "location": {
                        "lat": points2[i].latitude,
                        "lng": points2[i].longitude
                    },
                    "stopover": true
                }
                waypointsArrFinal.push(waypoint)
            }
            
            var request = {
                // from: Blackpool to: Preston to: Blackburn
                origin: {"lat": points[currentPosition].latitude,
                    "lng": points[currentPosition].longitude},
                destination: pos,
                waypoints: waypointsArrFinal,
                optimizeWaypoints: true,
                travelMode: google.maps.DirectionsTravelMode.WALKING
            }

            index++

            var service = new google.maps.DirectionsService()
            service.route(request, function (response, status) {
                if (status == google.maps.DirectionsStatus.OK) {
                    //  directionsDisplay.setDirection(response)
                    var route = response.routes[0]
                    var data = computeTotalDistance(response)

                    Materialize.toast('Success! On to the next point.', 4000)

                    var total_dist = 0
                    var total_time = 0
                    console.log(data)
                    displayNext(points[index], index * 10, (data[1]/1000) + " ", Math.round(Number(data[0]/60)) + " ")

                    currentPosition++
                }
            })
        } else {
            Materialize.toast('Not quite! Keep looking.', 4000)
        }
    }).fail(function (jqXHr, textStatus, errorThrown) {
        console.log('fail')
    }).always(function () { })
}