document.getElementById("nameChangeForm").addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent default form submission

    // Get form data
    const formData = {
        id: document.getElementById("name").value,                    
    };
    
    // Convert form data to JSON
    const jsonData = JSON.stringify(formData);

    // Send POST request
    fetch("/name_change/", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: jsonData
    })
    .then(response => {        

        if (response.redirected) {
            const redirectUrl = response.url;
            window.location.href = redirectUrl;
            } 
        else {
            return response.json();
        }
    })
    .catch(error => {      
    console.error("Error:", error);
    });
});

document.getElementById("notifyChangeForm").addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent default form submission

    // Get form data
    const formData = {
        transaction: document.getElementById("transaction").checked,
        daily: document.getElementById("daily").checked,
    };
    
    // Convert form data to JSON
    const jsonData = JSON.stringify(formData);

    // Send POST request
    fetch("/notify_change/", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: jsonData
    })
    .then(response => {        

        if (response.redirected) {
            const redirectUrl = response.url;
            window.location.href = redirectUrl;
            } 
        else {
            return response.json();
        }
    })
    .catch(error => {      
    console.error("Error:", error);
    });
});