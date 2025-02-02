// pages/invoiceJobAdd.js
document.addEventListener('DOMContentLoaded', function() {
    document.getElementById('jobForm').addEventListener('submit', function(event) {
        event.preventDefault();
        
        const jobName = document.getElementById('jobName').value;
        const price = document.getElementById('price').value;
        const status = document.getElementById('status').value === '1';
        
        fetch('/job-add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ jobName, price, status }),
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                window.location.href = '/pages/invoiceJob.html'; // Redirect to job list
            } else {
                alert('Failed to add job.');
            }
        })
        .catch(error => {
            console.error('Error:', error);
        });
    });
});
