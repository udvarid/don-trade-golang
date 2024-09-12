document.getElementById("nameChangeForm").addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent default form submission

    // Get form data
    const formData = {
        id: document.getElementById("id").value,                    
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