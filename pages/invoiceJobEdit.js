// pages/invoiceJobEdit.js
// Fetch job details when page loads
document.addEventListener('DOMContentLoaded', function() {
    const urlParams = new URLSearchParams(window.location.search);
    const jobName = urlParams.get('jobName'); // Use jobName instead of id

    if (jobName) {
        fetch(`/jobs`)
            .then(response => response.json())
            .then(jobs => {
                const job = jobs.find(j => j.jobName === jobName); // Find the job by jobName
                if (job) {
                    document.getElementById('jobId').value = job.id;
                    document.getElementById('jobName').value = job.jobName;
                    document.getElementById('price').value = job.price;
                } else {
                    console.error('Job not found');
                }
            })
            .catch(error => console.error('Error fetching jobs:', error));
    }
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
