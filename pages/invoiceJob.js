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
            <td><button onclick="window.location.href='/pages/invoice-job-edit?id=${job.id}'">Edit</button></td>
        `;
        tbody.appendChild(row);
    });
})
.catch(error => console.error('Error fetching jobs:', error));