// pages/invoiceJob.js
// Fetch and display jobs
fetch('/jobs')
.then(response => response.json())
.then(data => {
    const tbody = document.querySelector('#jobTable tbody');
    tbody.innerHTML = ''; // Clear existing rows
    data.forEach(job => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${job.jobName}</td>
            <td>${job.price}</td>
            <td><button onclick="window.location.href='/pages/invoice-job-edit?jobName=${encodeURIComponent(job.jobName)}'">Edit</button></td>
        `;
        tbody.appendChild(row);
    });
})
.catch(error => console.error('Error fetching jobs:', error));

function exportJobs() {
    fetch('/job-export')
        .then(response => response.blob())
        .then(blob => {
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'jobs.csv';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            window.URL.revokeObjectURL(url);
        })
        .catch(error => console.error('Error exporting jobs:', error));
}
