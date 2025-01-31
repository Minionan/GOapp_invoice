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

// Jobs exporing to CSV
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

// Jobs imporing from CSV
function importJobs() {
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = '.csv';
    fileInput.onchange = (event) => {
        const file = event.target.files[0];
        if (!file) return;

        const formData = new FormData();
        formData.append('file', file);

        fetch('/job-import', {
            method: 'POST',
            body: formData,
        })
        .then(response => response.json())
        .then(data => {
            let message = "Import Results:\n";
            message += `Imported Jobs: ${(data.imported || []).join(', ') || 'None'}\n`;
            message += `Skipped Jobs (duplicates): ${(data.skipped || []).join(', ') || 'None'}`;
            alert(message);
            window.location.href = '/pages/invoiceJob.html'; // Reload invoiceJob page
        })
        .catch(error => console.error('Error importing jobs:', error));
    };

    fileInput.click();
}
