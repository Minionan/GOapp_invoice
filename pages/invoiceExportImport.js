// pages/invoiceExpImp.js
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
                    <td>${invoice.parentName}</td>
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

function exportInvoices() {
    // Trigger the export endpoint
    window.location.href = '/invoice-export-csv';
}

function importInvoices() {
    const fileInput = document.createElement('input');
    fileInput.type = 'file';
    fileInput.accept = '.csv';
    fileInput.onchange = (event) => {
        const file = event.target.files[0];
        if (file) {
            const formData = new FormData();
            formData.append('file', file);

            $.ajax({
                url: '/invoice-import-csv',
                type: 'POST',
                data: formData,
                processData: false,
                contentType: false,
                success: function(response) {
                    let message = response.message;
                    if (response.importedInvoices && response.importedInvoices.length > 0) {
                        message += `\nImported invoices:\n${response.importedInvoices.join('\n')}`;
                    }
                    if (response.duplicateInvoices && response.duplicateInvoices.length > 0) {
                        message += `\nDuplicate invoices:\n${response.duplicateInvoices.join('\n')}`;
                    }
                    if (response.malformedRows && response.malformedRows.length > 0) {
                        message += `\nMalformed rows at lines:\n${response.malformedRows.join('\n')}`;
                    }
                    alert(message);
                    location.reload();
                },
                error: function() {
                    alert('Failed to import invoices. Please check the file format.');
                }
            });
        }
    };
    fileInput.click();
}
