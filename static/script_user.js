function handleChangeChartButtonClick(item) {
    const divs = document.querySelectorAll(`#barChart div`);        
    divs.forEach(div => {                
        if (div.id.startsWith("chart-")) {
            if ("chart-" + item === div.id) {            
                div.style.display = 'block';
            } else {                            
                div.style.display = 'none';
            }            
        }   
    });    
}


function handleClearButtonClick(pureItem, short) {
    console.log("Clear button clicked with parameter:", pureItem);
    
    fetch("/clear_item/"+pureItem+"/"+short, {
        method: "GET",
        headers: {
            "Content-Type": "application/json"
        }
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
}


function handleModifyOrderButtonClick(id) {
    const limitPriceElement = document.getElementById("limitPrice-"+id);
    const limitPrice = limitPriceElement ? limitPriceElement.value : "0";
    const validDays = document.getElementById("validDays-"+id).value;

    const formData = {
        order_id: id,
        limit_price: limitPrice,
        valid_days: validDays
    };
    
    // Convert form data to JSON
    const jsonData = JSON.stringify(formData); 
    
    fetch("/modify_order", {
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
}


document.addEventListener("DOMContentLoaded", function() {
    const marketType = document.getElementById("marketType");
    const orderType = document.getElementById("orderType");
    const limitPrice = document.getElementById("limitPrice");
    const orderForm = document.getElementById("orderForm");
    const usd = document.getElementById("usd");
    const allin = document.getElementById("allin");
    const short = document.getElementById("short");
    const number = document.getElementById("numberOfAssets");
    
    function toggleLimitPrice() {
        if (marketType.value === "MARKET") {
            limitPrice.disabled = true;
            limitPrice.value = "";
        } else {
            limitPrice.disabled = false;
        }
    }

    function toggleAllin() {
        if (allin.checked === true) {
            usd.value = "";
            usd.disabled = true;
            if (orderType.value === "SELL") {
                number.value = "";
                number.disabled = true;
            }
        } else {
            if (orderType.value === "BUY") {
                usd.disabled = false;                
            } else {
                number.disabled = false;
            }
        }
    }

    function toggleOrderType() {
        allin.disabled = false;
        if (orderType.value === "SELL") {
            usd.value = "";
            usd.disabled = true;
            short.disabled = false;          
        } else {
            if (allin.checked === false) {
                usd.disabled = false;
            }
            short.disabled = true;
            short.checked = false;
        }
    }

    function toggleShort() {
        if (short.checked === true) {
            allin.checked = false;
            allin.disabled = true;
            number.disabled = false;
        } else {
            allin.disabled = false;
        }
    }

    marketType.addEventListener("change", toggleLimitPrice);
    allin.addEventListener("change", toggleAllin);
    orderType.addEventListener("change", toggleOrderType);
    short.addEventListener("change", toggleShort);

    orderForm.addEventListener("submit", function(event) {        
        event.preventDefault();
        if (allin.checked === false && usd.value === "" && number.value === "") {
            alert("The number or the usd should be set or the all-in should be checked!");            
        } else {
            const formData = {
                item: document.getElementById("assetList").value,
                direction: document.getElementById("orderType").value,
                type: document.getElementById("marketType").value,
                limit_price: document.getElementById("limitPrice").value,
                number_of_items: document.getElementById("numberOfAssets").value,
                usd: document.getElementById("usd").value,
                all_in: document.getElementById("allin").checked,
                short: document.getElementById("short").checked,
                valid_days: document.getElementById("validDays").value,
            };
            
            // Convert form data to JSON
            const jsonData = JSON.stringify(formData); 
            
            fetch("/addorder", {
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
        }         
    });
    
    toggleLimitPrice();
    toggleAllin();
    toggleOrderType();
});
