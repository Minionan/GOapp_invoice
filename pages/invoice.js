// pages/invoice.js
$(document).ready(function() {
    // Fetch invoices from the server
    $.getJSON('/invoice-list', function(data) {
        const tableBody = $('#invoice-table tbody');
        tableBody.empty(); // Clear any existing rows

        // Populate the table with invoice data
        data.forEach(function(invoice) {
            const row = `
                <tr>
                    <td>${invoice.invoiceNumber}</td>
                    <td>${invoice.clientName}</td>
                    <td>${invoice.parentName}</td>
                    <td>${invoice.phone}</td>
                    <td>${invoice.email}</td>
                    <td>${invoice.cost}</td>
                    <td>${invoice.total}</td>
                </tr>
            `;
            tableBody.append(row);
        });
    }).fail(function() {
        alert('Failed to load invoices. Please try again.');
    });
});
