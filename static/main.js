const url = "portal.altbox.online:3000"
const endpoint = "/v1/recommender"

// click handler for recommend button
function recommendMovie() {
    movie = $("#txt-movie-input").val()
    console.log("recommend movie button clicked with movie = ", movie)
    sendRequest(url, movie)
}

// send request to backend server
function sendRequest(url, movie) {
    payload = {
        "name": movie
    }

    $.ajax(
        {
            type: "POST",
            url: url + endpoint,
            contentType: "application/json",
            data: JSON.stringify(payload),
            success: function(data) {
                console.log(data)
            }
        })
}
