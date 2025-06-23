/*
const addData = async (event) => {
    event.preventDefault();

    const AssetId = document.getElementById("AssetId").value;
    const GSTIN = document.getElementById("GSTIN").value;
    const Quantity = document.getElementById("Quantity").value;
    const InvoiceDate = document.getElementById("InvoiceDate").value;
    const Weight = document.getElementById("Weight").value;
    const Status = document.getElementById("Status").value;


    const ProductData = {
        AssetId: AssetId,
        GSTIN: GSTIN,
        Quantity: Quantity,
        InvoiceDate: InvoiceDate,
        Weight: Weight,
        Status: Status,
    };


    if (
        AssetId.length == 0 ||
        GSTIN.length == 0 ||
        Quantity.length == 0 ||
        InvoiceDate.length == 0 ||
        Weight.length == 0 ||
        Status.length == 0
    ) {
        alert("Please enter the data properly.");
    } else {
        try {
            const response = await fetch("/api/prada", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(ProductData),
            });
            console.log("RESPONSE: ", response)
            const data = await response.json()
            console.log("DATA: ", data);
            return alert("Car Created");
        } catch (err) {
            alert("Error");
            console.log(err);
        }
    }
};

const readData = async (event) => {
    event.preventDefault();
    const AssetId = document.getElementById("carIdInput").value;

    if (AssetId.length == 0) {
        alert("Please enter a valid ID.");
    } else {
        try {
            const response = await fetch(`/api/prada/${AssetId}`);
            let responseData = await response.json();
            console.log("sadsdsdd", responseData);
            alert(JSON.stringify(responseData));
        } catch (err) {
            alert("Error");
            console.log(err);
        }
    }
};
*/
////////

// function toManuDash() {
//     window.location.href = '/Status';
// }

function swalBasic(data) {
    swal.fire({
        // toast: true,
        icon: `${data.icon}`,
        title: `${data.title}`,
        animation: true,
        position: 'center',
        showConfirmButton: true,
        footer: `${data.footer}`,
       
    });
}

function swalBasicRefresh(data) {
    swal.fire({
        // toast: true,
        icon: `${data.icon}`,
        title: `${data.title}`,
        animation: true,
        position: 'center',
        showConfirmButton: true,
        footer: `${data.footer}`,
    }).then(() => {
        location.reload();
    });
}

function reloadWindow() {
    window.location.reload();
}

const addPradaData = async (event) => {
    event.preventDefault();
    const AssetId = document.getElementById('AssetId').value;
    const GSTIN = document.getElementById('GSTIN').value;
    const Quantity = document.getElementById('Quantity').value;
    const InvoiceDate = document.getElementById('InvoiceDate').value;
    const Weight = document.getElementById('Weight').value;
    const Status = document.getElementById('Status').value;
    const Destination = document.getElementById('Destination').value;
    console.log(AssetId + GSTIN + Quantity + InvoiceDate+Weight+Status+Destination);

    const Productdata = {
        AssetId: AssetId,
        GSTIN: GSTIN,
        Quantity: Quantity,
        InvoiceDate: InvoiceDate,
        Weight: Weight,
        Status:Status,
        Destination:Destination
    };
    if (
        AssetId.length == 0 ||
        GSTIN.length == 0 ||
        Quantity.length == 0 ||
        InvoiceDate.length == 0 ||
        Weight.length == 0 ||
        Status.length == 0||
        Destination.length ==0
    ) {
        const data = {
            title: "You might have missed something",
            footer: "Enter all mandatory fields to add a Prada",
            icon: "warning"
        }
        swalBasic(data);
    } else {
        try {
            const response = await fetch("/api/prada", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(Productdata),

            });
            console.log("RESPONSE: ", response)
            console.log("-----------==========----------")
            const data = await response.json()
            console.log("DATA: ", data);
            const pradaStatus = {
                title: "Success",
                footer: "Added a new prada",
                icon: "success"
            }
            swalBasicRefresh(pradaStatus);

        } catch (err) {
            // alert("Error");
            console.log(err);
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
        }
    }

}

const readProductData = async (event) => {
    event.preventDefault();
    const AssetId = document.getElementById("queryPradaId").value;
    if (AssetId.length == 0) {
        const data = {
            title: "Enter a valid prada Id",
            footer: "This is a mandatory field",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        try {
            const response = await fetch(`/api/prada/${AssetId}`);
            let responseData = await response.json();
            console.log("response", responseData);
            // alert(JSON.stringify(responseData));
            const dataBuf = JSON.stringify(responseData)
            swal.fire({
                // toast: true,
                icon: `success`,
                title: `Current status of prada with AssetId ${AssetId} :`,
                animation: false,
                position: 'center',
                html: `<h3>${dataBuf}</h3>`,
                showConfirmButton: true,
            })
        } catch (err) {

            console.log(err);
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
        }
    }
};

function getMatchingOrders(AssetId) {
    // console.log("AssetId", AssetId)
    window.location.href = '/api/order/match-prada?AssetId=' + AssetId;
}

//Method to get the history of an item
function getCarHistory(AssetId) {
    console.log("AssetId====", AssetId)
    window.location.href = '/api/prada/history?AssetId=' + AssetId;
}


const registerCar = async (event) => {
    // function registerCar(event) {
    console.log("Entered the register function")
    event.preventDefault();
    const AssetId = document.getElementById('AssetId').value;
    const carOwner = document.getElementById('carOwner').value;
    const regNumber = document.getElementById('regNumber').value;
    console.log(AssetId + carOwner + regNumber);
    const ProductData = {
        AssetId: AssetId,
        carOwner: carOwner,
        regNumber: regNumber,
    };
    if (AssetId.length == 0 || carOwner.length == 0 || regNumber.length == 0) {
        const data = {
            title: "You have missed something",
            footer: "All fields are mandatory",
            icon: "warning"
        }
        swalBasic(data)
    }
    else {
        try {
            const response = await fetch("/api/prada/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(ProductData),
            });
            console.log("RESPONSE: ", response)
            const data = await response.json()
            console.log("DATA: ", data);
            const pradaStatus = {
                title: `Registered prada ${AssetId} to ${carOwner}`,
                footer: "Registered prada",
                icon: "success"
            }
            swalBasicRefresh(pradaStatus);

        } catch (err) {
            console.log(err);
            const data = {
                title: `Failed to register prada`,
                footer: "Please try again !!",
                icon: "error"
            }
            swalBasic(data);
        }


    }
}


const addOrder = async (event) => {
    event.preventDefault();
    const OrderID = document.getElementById('OrderID').value;
    const Quantity = document.getElementById('Quantity').value;
    const Status = document.getElementById('Status').value;
    const DateOfOrder = document.getElementById('DateOfOrder').value;
    const DealerName = document.getElementById('DealerName').value;
    console.log(OrderID + Quantity + Status+DateOfOrder+DealerName);

    const orderData = {
        OrderID: OrderID,
        Quantity: Quantity,
        Status: Status,
        DateOfOrder: DateOfOrder,
        DealerName: DealerName,
    };
    if (
        OrderID.length == 0 ||
        Quantity.length == 0 ||
        Status.length == 0 ||
        DateOfOrder.length == 0 ||
        DealerName.length == 0
    ) {
        const data = {
            title: "You have missed something",
            footer: "All fields are mandatory",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        try {
            const response = await fetch("/api/order", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(orderData),
            });
            console.log("RESPONSE: ", response)
            const data = await response.json()
            console.log("DATA: ", data);
            // return alert("Order Created");
            const orderStatus = {
                title: `Order is created`,
                footer: "Raised Order",
                icon: "success"
            }
            swalBasic(orderStatus)
        } catch (err) {
            // alert("Error");
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
            console.log(err);
        }
    }

}

const readOrder = async (event) => {
    event.preventDefault();
    const orderId = document.getElementById("ordNum").value;

    if (orderId.length == 0) {
        const data = {
            title: "Enter a order number",
            footer: "Order Number is mandatory",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        try {
            const response = await fetch(`/api/order/${orderId}`);
            let responseData = await response.json();
            console.log("response", responseData);
            const dataBuf = JSON.stringify(responseData)
            swal.fire({
                // toast: true,
                icon: `success`,
                title: `Current status of Order : `,
                animation: false,
                position: 'center',
                html: `<h3>${dataBuf}</h3>`,
                showConfirmButton: true,
                timer: 3000,
                timerProgressBar: true,
                didOpen: (toast) => {
                    toast.addEventListener('mouseenter', swal.stopTimer)
                    toast.addEventListener('mouseleave', swal.resumeTimer)
                }
            })
            // alert(JSON.stringify(responseData));
        } catch (err) {
            // alert("Error");
            const data = {
                title: "Error in processing Request",
                footer: "Something went wrong !",
                icon: "error"
            }
            swalBasic(data);
            console.log(err);
        }
    }
};


async function matchOrder(orderId, AssetId) {
    if (!orderId || !AssetId) {
        const data = {
            title: "Enter a order number",
            footer: "Order Number is mandatory",
            icon: "warning"
        }
        swalBasic(data)
    } else {
        const matchData = {
            AssetId: AssetId,
            orderId: orderId,
        }
        try {
            const response = await fetch("/api/prada/match-order", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(matchData),
            });
            const data = {
                title: `Order matched successfully`,
                footer: "Order matched",
                icon: "success"
            }
            swalBasicRefresh(data)

        } catch (err) {
            const data = {
                title: `Failed to match order`,
                footer: "Please try again !!",
                icon: "error"
            }
            swalBasic(data)
        }


    }
}


function allOrders() {
    window.location.href = '/api/order/all';
}


async function getEvent() {
    try {
        const response = await fetch("/api/event", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            }
        });
        // console.log("RESPONSE: ", response)
        const data = await response.json()
        // console.log("DATA: ", data);

        const eventsData = data["carEvent"]
        swal.fire({
            toast: true,
            // icon: `${data.icon}`,
            title: `Event : `,
            animation: false,
            position: 'top-right',
            html: `<h5>${eventsData}</h5>`,
            showConfirmButton: false,
            timer: 5000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', swal.stopTimer)
                toast.addEventListener('mouseleave', swal.resumeTimer)
            }
        })
    } catch (err) {
        swal.fire({
            toast: true,
            icon: `error`,
            title: `Error`,
            animation: false,
            position: 'top-right',
            showConfirmButton: true,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', swal.stopTimer)
                toast.addEventListener('mouseleave', swal.resumeTimer)
            }
        })
        console.log(err);
    }
}





