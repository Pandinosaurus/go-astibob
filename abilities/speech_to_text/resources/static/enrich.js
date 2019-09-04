let enrich = {

    init: function() {
        base.init(enrich.onMessage, function() {
            // Get references
            asticode.tools.sendHttp({
                method: "GET",
                url: "../routes/references/enrich",
                error: base.httpError,
                success: function(data) {
                    // Update options
                    enrich.options = data.responseJSON

                    // Update store samples
                    let e = document.querySelector("#store-samples")
                    e.className = "toggle " + (enrich.options.store_samples ? "on": "off")
                    e.addEventListener("click", enrich.handleStoreSampleOption)

                    // Finish
                    base.finish()
                }
            })
        })
    },
    handleStoreSampleOption: function() {
        asticode.tools.sendHttp({
            method: "PATCH",
            url: "../routes/options/enrich",
            payload: JSON.stringify({
                store_samples: !enrich.options.store_samples,
            }),
            error: base.httpError,
        })
    },
    onMessage: function(data) {
        switch (data.name) {
            case "speech_to_text.enrich.options.updated":
                enrich.options = data.payload
                document.querySelector("#store-samples").className = "toggle " + (enrich.options.store_samples ? "on": "off")
                break
        }
    },
}