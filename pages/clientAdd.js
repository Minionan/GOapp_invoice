// pages/clientAdd.js
// Check if abbreviation is at least 3 characters long
function validateClientAbbreviation(abbreviation) {
    if (abbreviation.length < 3) {
        return 'Abbreviation must be at least 3 characters long.';
    }

    // Check if abbreviation contains only allowed characters (alphanumeric and underscores)
    const allowedCharacters = /^[a-zA-Z0-9_]+$/;
    if (!allowedCharacters.test(abbreviation)) {
        return 'Abbreviation can only contain letters, numbers, and underscores.';
    }

    return null; // No error
}

document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('clientForm').addEventListener('submit', function(event) {
        event.preventDefault();
        
        const clientName = document.getElementById('clientName').value;
        const parentName = document.getElementById('parentName').value;
        const address1 = document.getElementById('address1').value;
        const address2 = document.getElementById('address2').value;
        const phone = document.getElementById('phone').value;
        const email = document.getElementById('email').value;
        const abbreviation = document.getElementById('abbreviation').value;
        const status = document.getElementById('status').value === '1';

        // Validate abbreviation
        const abbreviationError = validateClientAbbreviation(abbreviation);
        if (abbreviationError) {
            alert(abbreviationError);
            return;
        }

        fetch('/client-add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ clientName, parentName, address1, address2, phone, email, abbreviation, status }),
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                window.location.href = '/pages/client.html'; // Redirect to client list
            } else {
                alert('Failed to add client: ' + (data.error || 'Unknown error'));
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
    });
});
