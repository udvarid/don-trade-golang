function showDiv(selectElement, pagesArray) {
    const selectedValue = selectElement.value;
    const divs = document.querySelectorAll(`#${selectElement.dataset.target} div`);
    pagesArray.forEach(page => {
        divs.forEach(div => {                    
            if (page.Name === div.id) {
                if (div.id === selectedValue) {                            
                    div.style.display = 'block';
                } else {                            
                    div.style.display = 'none';
                }
            }
        });
    });
}

document.getElementById("loginForm").addEventListener("submit", function(event) {
    event.preventDefault(); // Prevent default form submission

    // Get form data
    const formData = {
        id: document.getElementById("id").value,                    
    };
    
    // Convert form data to JSON
    const jsonData = JSON.stringify(formData);

    // Send POST request
    fetch("/validate/", {
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
        } else {
            return response.json();
        }
    })
    .catch(error => {      
    console.error("Error:", error);
    });
});