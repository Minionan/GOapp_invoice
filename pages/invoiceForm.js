// pages/invoiceForm.js
$(document).ready(function() {
    let maxJobRows = 1; // Default value

    // Load jobs data
    $.getJSON('/jobs', function(data) {
        window.jobsData = data;
        updateJobDropdowns();
    });

    // Load clients data
    $.getJSON('/clients', function(data) {
        window.clientsData = data;
        var clientDropdown = $('#client-dropdown');
        data.forEach(function(client) {
            clientDropdown.append($('<option>', {
                value: client.abbreviation,
                text: client.clientName
            }));
        });
    });

    // Load maximum number of rows from template.xlsx
    $.getJSON('/max-job-rows', function(data) {
        maxJobRows = data.maxRows;
    });

    function updateJobDropdowns() {
        $('.job-dropdown').each(function() {
            var dropdown = $(this);
            // Only update if the dropdown is empty or has no selection
            if (!dropdown.val()) {
                dropdown.empty();
                dropdown.append($('<option>', {
                    value: "",
                    text: "Please choose a job..."
                }));
                if (window.jobsData) {
                    window.jobsData.forEach(function(job) {
                        dropdown.append($('<option>', {
                            value: job.jobName,
                            text: job.jobName,
                            'data-price': job.price
                        }));
                    });
                }
            }
        });
    }

    function updateClientDetailsAndInvoiceNumber() {
        let selectedClient = $('#client-dropdown').val();
        let selectedDate = $('#invoiceDate').val();
        
        // Clear fields if no client selected
        if (!selectedClient || !window.clientsData) {
            $('#client-details').html('');
            $('#invoiceNumber').val('');
            return;
        }

        // Update client details
        let client = window.clientsData.find(c => c.abbreviation === selectedClient);
        if (client) {
            $('#client-details').html(
                '<strong>Payee Name:</strong> ' + client.parentName + '<br>' +
                '<strong>Address 1:</strong> ' + client.address1 + '<br>' +
                (client.address2 ? '<strong>Address 2:</strong> ' + client.address2 + '<br>' : '') +
                '<strong>Phone:</strong> ' + client.phone + '<br>' +
                '<strong>Email:</strong> ' + client.email
            );
            
            // Update invoice number if date is selected
            if (selectedDate) {
                let dateObj = new Date(selectedDate);
                let year = dateObj.getFullYear();
                let month = String(dateObj.getMonth() + 1).padStart(2, '0');
                $('#invoiceNumber').val(`${year}_${month}_${client.abbreviation}`);
            } else {
                $('#invoiceNumber').val('');
            }
        }
    }

    // Event Handlers
    $('#client-dropdown, #invoiceDate').change(function() {
        updateClientDetailsAndInvoiceNumber();
    });

    $('#job-rows').on('change', '.job-dropdown', function() {
        var row = $(this).closest('.job-row');
        var selectedJob = window.jobsData.find(job => job.jobName === $(this).val());
        if (selectedJob) {
            row.find('.price').val(selectedJob.price);
            updateFullPrice(row);
        }
    });

    $('#add-job').click(function() {
        if ($('.job-row').length >= maxJobRows) {
            alert(`Maximum number of jobs ${maxJobRows} has been reached for this invoice form.`);
            return;
        }
        var jobRow = $('.job-row:first').clone();
        jobRow.find('.quantity, .price, .fullPrice').val('');
        jobRow.append('<button type="button" class="delete-row" style="margin-left: 10px;">Delete</button>');
        $('#job-rows').append(jobRow);
        updateJobDropdowns();
    });

    $('#job-rows').on('click', '.delete-row', function() {
        if ($('.job-row').length > 1) {
            $(this).closest('.job-row').remove();
            calculateTotals();
        }
    });

    $('#job-rows').on('input', '.quantity, .price', function() {
        var row = $(this).closest('.job-row');
        updateFullPrice(row);
    });

    function calculateTotals() {
        var totalCost = 0;
        $('.job-row').each(function() {
            totalCost += parseFloat($(this).find('.fullPrice').val()) || 0;
        });

        var vat = totalCost * 0.05;
        var totalAmount = totalCost + vat;

        $('#cost').val(totalCost.toFixed(2));
        $('#vat').val(vat.toFixed(2));
        $('#total').val(totalAmount.toFixed(2));
    }

    function updateFullPrice(row) {
        var quantity = parseFloat(row.find('.quantity').val()) || 0;
        var price = parseFloat(row.find('.price').val()) || 0;
        var fullPrice = quantity * price;
        row.find('.fullPrice').val(fullPrice.toFixed(2));
        calculateTotals();
    }

    function updateClientDetailsAndInvoiceNumber() {
        let selectedClient = $('#client-dropdown').val();
        let selectedDate = $('#invoiceDate').val();
        
        // Clear fields if no client selected
        if (!selectedClient || !window.clientsData) {
            $('#client-details').html('');
            $('#invoiceNumber').val('');
            return;
        }
    
        // Update client details
        let client = window.clientsData.find(c => c.abbreviation === selectedClient);
        if (client) {
            $('#client-details').html(
                '<strong>Parent Name:</strong> ' + client.parentName + '<br>' +
                '<strong>Address 1:</strong> ' + client.address1 + '<br>' +
                (client.address2 ? '<strong>Address 2:</strong> ' + client.address2 + '<br>' : '') +
                '<strong>Phone:</strong> ' + client.phone + '<br>' +
                '<strong>Email:</strong> ' + client.email
            );
            
            // Update invoice number if date is selected
            if (selectedDate) {
                let dateObj = new Date(selectedDate);
                let year = dateObj.getFullYear();
                let month = String(dateObj.getMonth() + 1).padStart(2, '0');
                let invoiceNumber = `${year}_${month}_${client.abbreviation}`;
                $('#invoiceNumber').val(invoiceNumber);
    
                // Check if the invoice number already exists
                $.getJSON('/invoice-number-check?invoiceNumber=' + invoiceNumber, function(data) {
                    if (data.exists) {
                        alert('An invoice with the same number already exists. Please choose a different date or client. If you want to change the existing invoice, please use the Invoice Edit Form.');
                        $('#client-dropdown').val(''); // Clear client dropdown
                        $('#invoiceDate').val(''); // Clear date field
                        $('#invoiceNumber').val(''); // Clear invoice number field
                        $('#client-details').html(''); // Clear client details
                    }
                });
            } else {
                $('#invoiceNumber').val('');
            }
        }
    }

    $('#generate-txt').click(function() {
        // Check if invoice number is set
        if (!$('#invoiceNumber').val()) {
            alert('Please select both client and date to generate invoice number before downloading the file.');
            return;
        }

        // Check for empty or zero values in job rows
        let issues = [];
        let hasIssues = false;

        $('.job-row').each(function(index) {
            const rowNum = index + 1;
            const jobDescription = $(this).find('.job-dropdown').val();
            const quantity = parseFloat($(this).find('.quantity').val()) || 0;
            const price = parseFloat($(this).find('.price').val()) || 0;
            const fullPrice = parseFloat($(this).find('.fullPrice').val()) || 0;

            if (!jobDescription) {
                issues.push(`Row ${rowNum}: Job Description is not selected`);
                hasIssues = true;
            }
            if (quantity === 0) {
                issues.push(`Row ${rowNum}: Quantity is zero`);
                hasIssues = true;
            }
            if (price === 0) {
                issues.push(`Row ${rowNum}: Price is zero`);
                hasIssues = true;
            }
            if (fullPrice === 0) {
                issues.push(`Row ${rowNum}: Full Price is zero`);
                hasIssues = true;
            }
        });

        if (hasIssues) {
            const message = 'The following issues were found:\n\n' + 
                            issues.join('\n') + 
                            '\n\nDo you want to proceed with generating the file anyway?';
            
            if (!confirm(message)) {
                return;
            }
        }

        // Fetch the selected client details
        let selectedClient = $('#client-dropdown').val();
        let client = window.clientsData.find(c => c.abbreviation === selectedClient);

        // Proceed with file generation if no issues or user confirmed
        let invoiceDetails = `Invoice Number: ${$('#invoiceNumber').val()}\n\n`;
        invoiceDetails += `Invoice Date: ${$('#invoiceDate').val()}\n\n`;
        invoiceDetails += `Client Name: ${client.clientName}\n`;
        invoiceDetails += `Parent Name: ${client.parentName}\n`;
        invoiceDetails += `Address 1: ${client.address1}\n`;
        invoiceDetails += `Address 2: ${client.address2}\n`;
        invoiceDetails += `Phone: ${client.phone}\n`;
        invoiceDetails += `Email: ${client.email}\n\n`;
        invoiceDetails += "Jobs:\n";

        $('.job-row').each(function(index) {
            const jobDescription = $(this).find('.job-dropdown').val();
            const quantity = $(this).find('.quantity').val();
            const price = $(this).find('.price').val();
            const fullPrice = $(this).find('.fullPrice').val();

            if (jobDescription) {
                invoiceDetails += `${index + 1}: ${jobDescription}, Quantity: ${quantity}, Price: ${price}, Full Price: ${fullPrice}\n`;
            }
        });

        invoiceDetails += `\nJob Cost: ${$('#cost').val()}\n`;
        invoiceDetails += `VAT (5%): ${$('#vat').val()}\n`;
        invoiceDetails += `Total Amount: ${$('#total').val()}\n`;

        // Use the formatted string as the data for the blob
        var blob = new Blob([invoiceDetails], {type: 'text/plain'});
        var url = window.URL.createObjectURL(blob);
        var a = document.createElement('a');
        a.href = url;
        a.download = $('#invoiceNumber').val() + '.txt';
        document.body.appendChild(a);
        a.click();
        a.remove();
        window.URL.revokeObjectURL(url);
    });

    $('#generate-xlsx').click(function() {
        // Check if invoice number is set
        if (!$('#invoiceNumber').val()) {
            alert('Please select both client and date to generate invoice number before downloading the file.');
            return;
        }

        // Check for empty or zero values (same validation as in txt generation)
        let issues = [];
        let hasIssues = false;

        $('.job-row').each(function(index) {
            const rowNum = index + 1;
            const jobDescription = $(this).find('.job-dropdown').val();
            const quantity = parseFloat($(this).find('.quantity').val()) || 0;
            const price = parseFloat($(this).find('.price').val()) || 0;
            const fullPrice = parseFloat($(this).find('.fullPrice').val()) || 0;

            if (!jobDescription) {
                issues.push(`Row ${rowNum}: Job Description is not selected`);
                hasIssues = true;
            }
            if (quantity === 0) {
                issues.push(`Row ${rowNum}: Quantity is zero`);
                hasIssues = true;
            }
            if (price === 0) {
                issues.push(`Row ${rowNum}: Price is zero`);
                hasIssues = true;
            }
            if (fullPrice === 0) {
                issues.push(`Row ${rowNum}: Full Price is zero`);
                hasIssues = true;
            }
        });

        if (hasIssues) {
            const message = 'The following issues were found:\n\n' + 
                            issues.join('\n') + 
                            '\n\nDo you want to proceed with generating the file anyway?';
            
            if (!confirm(message)) {
                return;
            }
        }

        // Get the selected client details
        let selectedClient = $('#client-dropdown').val();
        let client = window.clientsData.find(c => c.abbreviation === selectedClient);

        // Collect job details
        let jobs = [];
        $('.job-row').each(function(index) {
            jobs.push({
                jobName: $(this).find('.job-dropdown').val(),
                quantity: $(this).find('.quantity').val(),
                price: $(this).find('.price').val(),
                fullPrice: $(this).find('.fullPrice').val()
            });
        });

        // Prepare the data for the XLSX generation
        let invoiceData = {
            parentName: client.parentName,
            address1: client.address1,
            address2: client.address2,
            phone: client.phone,
            email: client.email,
            invoiceNumber: $('#invoiceNumber').val(),
            invoiceDate: $('#invoiceDate').val(),
            cost: parseFloat($('#cost').val()),
            vat: parseFloat($('#vat').val()),
            total: parseFloat($('#total').val()),
            jobs: jobs  // Add the jobs array to the invoice data
        };

        // Send request to generate XLSX
        fetch('/generate-xlsx', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(invoiceData)
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
            a.download = $('#invoiceNumber').val() + '.xlsx';
            document.body.appendChild(a);
            a.click();
            a.remove();
            window.URL.revokeObjectURL(url);
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to generate XLSX file');
        });
    });

    $('#generate-pdf').click(function() {
        // Check if invoice number is set
        if (!$('#invoiceNumber').val()) {
            alert('Please select both client and date to generate invoice number before downloading the file.');
            return;
        }

        // Validation and data collection (same as in XLSX generation)
        let issues = [];
        let hasIssues = false;

        $('.job-row').each(function(index) {
            const rowNum = index + 1;
            const jobDescription = $(this).find('.job-dropdown').val();
            const quantity = parseFloat($(this).find('.quantity').val()) || 0;
            const price = parseFloat($(this).find('.price').val()) || 0;
            const fullPrice = parseFloat($(this).find('.fullPrice').val()) || 0;

            if (!jobDescription) {
                issues.push(`Row ${rowNum}: Job Description is not selected`);
                hasIssues = true;
            }
            if (quantity === 0) {
                issues.push(`Row ${rowNum}: Quantity is zero`);
                hasIssues = true;
            }
            if (price === 0) {
                issues.push(`Row ${rowNum}: Price is zero`);
                hasIssues = true;
            }
            if (fullPrice === 0) {
                issues.push(`Row ${rowNum}: Full Price is zero`);
                hasIssues = true;
            }
        });

        if (hasIssues) {
            const message = 'The following issues were found:\n\n' + 
                            issues.join('\n') + 
                            '\n\nDo you want to proceed with generating the file anyway?';
            
            if (!confirm(message)) {
                return;
            }
        }

        // Get the selected client details
        let selectedClient = $('#client-dropdown').val();
        let client = window.clientsData.find(c => c.abbreviation === selectedClient);

        // Collect job details
        let jobs = [];
        $('.job-row').each(function(index) {
            jobs.push({
                jobName: $(this).find('.job-dropdown').val(),
                quantity: $(this).find('.quantity').val(),
                price: $(this).find('.price').val(),
                fullPrice: $(this).find('.fullPrice').val()
            });
        });

        // Prepare the data for the PDF generation
        let invoiceData = {
            parentName: client.parentName,
            address1: client.address1,
            address2: client.address2,
            phone: client.phone,
            email: client.email,
            invoiceNumber: $('#invoiceNumber').val(),
            invoiceDate: $('#invoiceDate').val(),
            cost: parseFloat($('#cost').val()),
            vat: parseFloat($('#vat').val()),
            total: parseFloat($('#total').val()),
            jobs: jobs
        };

        // Send request to generate PDF
        fetch('/generate-pdf', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(invoiceData)
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
            a.download = $('#invoiceNumber').val() + '.pdf';
            document.body.appendChild(a);
            a.click();
            a.remove();
            window.URL.revokeObjectURL(url);
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to generate PDF file');
        });
    });

    $('#save-invoice').click(function() {
        // Check if invoice number is set
        if (!$('#invoiceNumber').val()) {
            alert('Please select both client and date to generate invoice number before saving the invoice.');
            return;
        }

        // Validation and data collection (same as in XLSX generation)
        let issues = [];
        let hasIssues = false;

        $('.job-row').each(function(index) {
            const rowNum = index + 1;
            const jobDescription = $(this).find('.job-dropdown').val();
            const quantity = parseFloat($(this).find('.quantity').val()) || 0;
            const price = parseFloat($(this).find('.price').val()) || 0;
            const fullPrice = parseFloat($(this).find('.fullPrice').val()) || 0;

            if (!jobDescription) {
                issues.push(`Row ${rowNum}: Job Description is not selected`);
                hasIssues = true;
            }
            if (quantity === 0) {
                issues.push(`Row ${rowNum}: Quantity is zero`);
                hasIssues = true;
            }
            if (price === 0) {
                issues.push(`Row ${rowNum}: Price is zero`);
                hasIssues = true;
            }
            if (fullPrice === 0) {
                issues.push(`Row ${rowNum}: Full Price is zero`);
                hasIssues = true;
            }
        });

        if (hasIssues) {
            const message = 'The following issues were found:\n\n' + 
                            issues.join('\n') + 
                            '\n\nDo you want to proceed with generating the file anyway?';
            
            if (!confirm(message)) {
                return;
            }
        }
    
        // Get the selected client details
        let selectedClient = $('#client-dropdown').val();
        let client = window.clientsData.find(c => c.abbreviation === selectedClient);
    
        // Collect job details
        let jobs = [];
        $('.job-row').each(function(index) {
            jobs.push({
                jobName: $(this).find('.job-dropdown').val(),
                quantity: $(this).find('.quantity').val(),
                price: $(this).find('.price').val(),
                fullPrice: $(this).find('.fullPrice').val()
            });
        });
    
        // Prepare the data for saving to the database
        let invoiceData = {
            invoiceNumber: $('#invoiceNumber').val(),
            invoiceDate: $('#invoiceDate').val(),
            clientName: client.clientName,
            parentName: client.parentName,
            address1: client.address1,
            address2: client.address2,
            phone: client.phone,
            email: client.email,
            cost: parseFloat($('#cost').val()),
            VAT: parseFloat($('#vat').val()),
            total: parseFloat($('#total').val()),
            jobs: jobs  // Add the jobs array to the invoice data
        };
    
        // Send request to save the invoice to the database
        fetch('/invoice-save', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(invoiceData)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                alert('Invoice saved to database successfully!');
            } else {
                alert('Failed to save invoice to database.');
            }
        })
        .catch(error => {
            console.error('Error:', error);
            alert('Failed to save invoice to database.');
        });
    });
});
