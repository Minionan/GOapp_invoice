// pages/invocieVat.js
document.addEventListener('DOMContentLoaded', function() {
    // Fetch the current VAT rate and populate the form
    fetch('/vat-get')
        .then(response => {
            if (!response.ok) {
                throw new Error('VAT rate not found');
            }
            return response.json();
        })
        .then(data => {
            document.getElementById('vatRate').value = data.rate || 0; // Default to 0 if rate is undefined
        })
        .catch(error => {
            console.error('Error:', error);
            document.getElementById('vatRate').value = 0; // Default to 0 on error
        });

    document.getElementById('vatForm').addEventListener('submit', function(event) {
        event.preventDefault();
        const vatRate = document.getElementById('vatRate').value;

        fetch('/vat-update', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ rate: parseFloat(vatRate) }), // Ensure rate is a number
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => { throw new Error(text); });
            }
            return response.json();
        })
        .then(data => {
            alert('VAT rate updated successfully!');
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to update VAT rate: ' + error.message);
        });
    });
});
