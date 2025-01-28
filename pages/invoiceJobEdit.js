// pages/invoiceJobEdit.js
// Fetch job details when page loads
document.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    const jobId = urlParams.get('id');
    
    fetch(`/job-details?id=${jobId}`)
        .then(response => response.json())
        .then(job => {
            document.getElementById('jobId').value = job.id;
            document.getElementById('jobName').value = job.jobName;
            document.getElementById('price').value = job.price;
        })
        .catch(error => console.error('Error fetching job details:', error));
});

// Save job chnages to the database
function saveJob() {
    const jobId = document.getElementById('jobId').value;
    const jobName = document.getElementById('jobName').value;
    const price = document.getElementById('price').value;

    fetch(`/job-update?id=${jobId}`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ jobName, price })
    })
    .then(response => {
        if (response.ok) {
            window.location.href = '/pages/invoiceJob.html';
        } else {
            throw new Error('Failed to update job');
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Failed to update job');
    });
}

// delete job from the database
function deleteJob() {
    const jobId = document.getElementById('jobId').value;
    const confirmDelete = confirm('Are you sure you want to delete this job?');

    if (confirmDelete) {
        fetch(`/job-delete?id=${jobId}`, {
            method: 'POST'
        })
        .then(response => {
            if (response.ok) {
                window.location.href = '/pages/invoiceJob.html';
            } else {
                throw new Error('Failed to delete job');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to delete job');
        });
    }
}
