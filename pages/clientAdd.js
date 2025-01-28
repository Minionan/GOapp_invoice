// pages/clientAdd.js
document.getElementById('clientForm').addEventListener('submit', function(event) {
    event.preventDefault();
    
    const clientName = document.getElementById('clientName').value;
    const parentName = document.getElementById('parentName').value;
    const address1 = document.getElementById('address1').value;
    const address2 = document.getElementById('address2').value;
    const phone = document.getElementById('phone').value;
    const email = document.getElementById('email').value;
    const abbreviation = document.getElementById('abbreviation').value;
    
    fetch('/client-add', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ clientName, parentName, address1, address2, phone, email, abbreviation }),
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            window.location.href = '/pages/client.html'; // Redirect to client list
        } else {
            alert('Failed to add client.');
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});
