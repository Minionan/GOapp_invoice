// pages/invoiceEdit.js
$(document).ready(function() {
    let vatRate = 0.05; // Default VAT rate

    // Fetch the list of invoices for the dropdown
    $.getJSON('/invoice-list', function(data) {
        const select = $('#invoice-select');
        select.empty(); // Clear any existing options
        select.append('<option value="">-- Select an Invoice --</option>');

        // Populate the dropdown with invoice numbers
        data.forEach(function(invoice) {
            select.append(`<option value="${invoice.invoiceNumber}">${invoice.invoiceNumber}</option>`);
        });
    }).fail(function() {
        alert('Failed to load invoices. Please try again.');
    });

    // Handle dropdown change event
    $('#invoice-select').change(function() {
        const invoiceNumber = $(this).val();
        if (invoiceNumber) {
            // Fetch the selected invoice details
            $.getJSON(`/invoice-list?invoiceNumber=${invoiceNumber}`, function(invoice) {
                populateInvoiceDetails(invoice);
            }).fail(function() {
                alert('Failed to fetch invoice details.');
            });
        }
    });

    // Function to populate editable fields with invoice details
    function populateInvoiceDetails(invoice) {
        $('#invoiceNumber').val(invoice.invoiceNumber);
        $('#invoiceDate').val(invoice.invoiceDate);
        $('#clientName').val(invoice.clientName);
        $('#parentName').val(invoice.parentName);
        $('#address1').val(invoice.address1);
        $('#address2').val(invoice.address2);
        $('#phone').val(invoice.phone);
        $('#email').val(invoice.email);
        $('#cost').val(invoice.cost.toFixed(2));
        $('#vat').val(invoice.vat.toFixed(2));
        $('#total').val(invoice.total.toFixed(2));

        // Calculate VAT rate based on the loaded invoice
        vatRate = invoice.vat / invoice.cost;

        // Populate job rows
        const jobRows = $('#job-rows');
        jobRows.empty();
        invoice.jobs.forEach((job, index) => {
            const jobRow = `
                <div class="job-row">
                    <input type="text" class="jobName" value="${job.jobName}">
                    <input type="number" class="quantity" value="${job.quantity}">
                    <input type="number" class="price" value="${job.price}">
                    <input type="number" class="fullPrice" value="${job.fullPrice}" readonly>
                    <button type="button" class="delete-job">Delete</button>
                </div>
            `;
            jobRows.append(jobRow);
        });
    }

    // Event listener for changes to quantity or price fields
    $('#job-rows').on('input', '.quantity, .price', function() {
        const row = $(this).closest('.job-row');
        updateFullPrice(row);
    });

    // Function to update fullPrice for a specific job row
    function updateFullPrice(row) {
        const quantity = parseFloat(row.find('.quantity').val()) || 0;
        const price = parseFloat(row.find('.price').val()) || 0;
        const fullPrice = quantity * price;
        row.find('.fullPrice').val(fullPrice.toFixed(2));

        // Recalculate totals
        calculateTotals();
    }

    // Function to recalculate totals (cost, VAT, total)
    function calculateTotals() {
        let totalCost = 0;
        $('.job-row').each(function() {
            totalCost += parseFloat($(this).find('.fullPrice').val()) || 0;
        });

        const vat = totalCost * vatRate;
        const totalAmount = totalCost + vat;

        $('#cost').val(totalCost.toFixed(2));
        $('#vat').val(vat.toFixed(2));
        $('#total').val(totalAmount.toFixed(2));
    }

    // Event listener for deleting a job row
    $('#job-rows').on('click', '.delete-job', function() {
        // Remove the closest job row
        $(this).closest('.job-row').remove();

        // Recalculate totals after deleting a row
        calculateTotals();
    });

    // Function to create a new empty job row
    function createEmptyJobRow() {
        return `
            <div class="job-row">
                <input type="text" class="jobName" placeholder="Job Description">
                <input type="number" class="quantity" placeholder="Quantity" min="0">
                <input type="number" class="price" placeholder="Price" min="0">
                <input type="number" class="fullPrice" placeholder="Full Price" readonly>
                <button type="button" class="delete-job">Delete</button>
            </div>
        `;
    }

    // Event listener for adding a new job row
    $('#add-job').click(function() {
        // Append a new empty job row
        $('#job-rows').append(createEmptyJobRow());
    });


    // Handle Save Invoice button click
    $('#save-invoice').click(function() {
        const invoiceNumber = $('#invoice-select').val();
        if (!invoiceNumber) {
            alert('Please select an invoice first.');
            return;
        }

        // Collect updated invoice data
        const updatedInvoice = {
            invoiceNumber: $('#invoiceNumber').val(),
            invoiceDate: $('#invoiceDate').val(),
            clientName: $('#clientName').val(),
            parentName: $('#parentName').val(),
            address1: $('#address1').val(),
            address2: $('#address2').val(),
            phone: $('#phone').val(),
            email: $('#email').val(),
            cost: parseFloat($('#cost').val()),
            vat: parseFloat($('#vat').val()),
            total: parseFloat($('#total').val()),
            jobs: []
        };

        // Collect job data
        $('.job-row').each(function() {
            const job = {
                jobName: $(this).find('.jobName').val(),
                quantity: $(this).find('.quantity').val(),
                price: $(this).find('.price').val(),
                fullPrice: $(this).find('.fullPrice').val()
            };
            updatedInvoice.jobs.push(job);
        });

        // Send request to update the invoice
        fetch('/invoice-update', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(updatedInvoice)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                alert('Invoice updated successfully!');
            } else {
                alert('Failed to update invoice.');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to update invoice.');
        });
    });

    // Handle Delete Invoice button click
    $('#delete-invoice').click(function() {
        const invoiceNumber = $('#invoice-select').val();
        if (!invoiceNumber) {
            alert('Please select an invoice first.');
            return;
        }

        // Confirm deletion
        if (confirm('Are you sure you want to delete this invoice? This action cannot be undone.')) {
            // Send request to delete the invoice
            fetch(`/invoice-delete?invoiceNumber=${invoiceNumber}`, {
                method: 'DELETE'
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                if (data.success) {
                    alert('Invoice deleted successfully!');
                    window.location.href = '/'; // Redirect to the main page
                } else {
                    alert('Failed to delete invoice.');
                }
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Failed to delete invoice.');
            });
        }
    });

    // pages/invoiceEdit.js
    $('#generate-txt').click(function() {
        const invoiceNumber = $('#invoice-select').val();
        if (!invoiceNumber) {
            alert('Please select an invoice first.');
            return;
        }

        // Fetch the selected invoice details
        $.getJSON(`/invoice-list?invoiceNumber=${invoiceNumber}`, function(invoice) {
            // Construct the text content
            let invoiceDetails = `Invoice Number: ${invoice.invoiceNumber}\n\n`;
            invoiceDetails += `Invoice Date: ${invoice.invoiceDate}\n\n`;
            invoiceDetails += `Client Name: ${invoice.clientName}\n`;
            invoiceDetails += `Parent Name: ${invoice.parentName}\n`;
            invoiceDetails += `Address 1: ${invoice.address1}\n`;
            invoiceDetails += `Address 2: ${invoice.address2}\n`;
            invoiceDetails += `Phone: ${invoice.phone}\n`;
            invoiceDetails += `Email: ${invoice.email}\n\n`;
            invoiceDetails += "Jobs:\n";

            invoice.jobs.forEach((job, index) => {
                invoiceDetails += `${index + 1}: ${job.jobName}, Quantity: ${job.quantity}, Price: ${job.price}, Full Price: ${job.fullPrice}\n`;
            });

            invoiceDetails += `\nJob Cost: ${invoice.cost.toFixed(2)}\n`;
            invoiceDetails += `VAT (5%): ${invoice.vat.toFixed(2)}\n`;
            invoiceDetails += `Total Amount: ${invoice.total.toFixed(2)}\n`;

            // Create a text file with the invoice details
            const blob = new Blob([invoiceDetails], { type: 'text/plain' });
            const url = window.URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = `${invoiceNumber}.txt`;
            document.body.appendChild(a);
            a.click();
            a.remove();
            window.URL.revokeObjectURL(url);
        }).fail(function() {
            alert('Failed to fetch invoice details.');
        });
    });

    $('#generate-xlsx').click(function() {
        const invoiceNumber = $('#invoice-select').val();
        if (!invoiceNumber) {
            alert('Please select an invoice first.');
            return;
        }

        // Fetch the selected invoice details
        $.getJSON(`/invoice-list?invoiceNumber=${invoiceNumber}`, function(invoice) {
            // Send request to generate XLSX
            fetch('/generate-xlsx', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(invoice)
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.blob();
            })
            .then(blob => {
                // Create a download link and trigger it
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = `${invoiceNumber}.xlsx`;
                document.body.appendChild(a);
                a.click();
                a.remove();
                window.URL.revokeObjectURL(url);
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Failed to generate XLSX file');
            });
        }).fail(function() {
            alert('Failed to fetch invoice details.');
        });
    });

    $('#generate-pdf').click(function() {
        const invoiceNumber = $('#invoice-select').val();
        if (!invoiceNumber) {
            alert('Please select an invoice first.');
            return;
        }

        // Fetch the selected invoice details
        $.getJSON(`/invoice-list?invoiceNumber=${invoiceNumber}`, function(invoice) {
            // Send request to generate PDF
            fetch('/generate-pdf', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(invoice)
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.blob();
            })
            .then(blob => {
                // Create a download link and trigger it
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = `${invoiceNumber}.pdf`;
                document.body.appendChild(a);
                a.click();
                a.remove();
                window.URL.revokeObjectURL(url);
            })
            .catch(error => {
                console.error('Error:', error);
                alert('Failed to generate PDF file');
            });
        }).fail(function() {
            alert('Failed to fetch invoice details.');
        });
    });
});
