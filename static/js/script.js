// DOM Elements
document.addEventListener('DOMContentLoaded', function() {
    // Modal Elements
    const deleteModal = document.getElementById('delete-modal');
    const confirmDeleteBtn = document.getElementById('confirm-delete');
    const cancelDeleteBtn = document.getElementById('cancel-delete');
    const filterForm = document.getElementById('filter-form');
    const resetFiltersBtn = document.getElementById('reset-filters');
    
    let transactionIdToDelete = null;

    // Handle transaction deletion
    const deleteButtons = document.querySelectorAll('.delete-transaction');
    deleteButtons.forEach(button => {
        button.addEventListener('click', function() {
            transactionIdToDelete = this.getAttribute('data-id');
            openModal();
        });
    });

    // Confirm deletion
    confirmDeleteBtn.addEventListener('click', function() {
        if (transactionIdToDelete) {
            // Send DELETE request to the server
            fetch(`/transaction/${transactionIdToDelete}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json'
                }
            })
            .then(response => {
                if (response.ok) {
                    // Remove the row from the table or reload the page
                    window.location.reload();
                } else {
                    console.error('Failed to delete transaction');
                    alert('Failed to delete transaction. Please try again.');
                }
                closeModal();
            })
            .catch(error => {
                console.error('Error:', error);
                alert('An error occurred. Please try again.');
                closeModal();
            });
        }
    });

    // Cancel deletion
    cancelDeleteBtn.addEventListener('click', closeModal);

    // Close modal if clicked outside
    window.addEventListener('click', function(event) {
        if (event.target === deleteModal) {
            closeModal();
        }
    });

    // Filter form submission
    if (filterForm) {
        filterForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            // Get form values
            const category = document.getElementById('filter-category').value;
            const type = document.getElementById('filter-type').value;
            const startDate = document.getElementById('start-date').value;
            const endDate = document.getElementById('end-date').value;
            
            // Build query parameters
            const params = new URLSearchParams();
            if (category) params.append('category', category);
            if (type) params.append('type', type);
            if (startDate) params.append('start_date', startDate);
            if (endDate) params.append('end_date', endDate);
            
            // Redirect to the same page with filters applied
            window.location.href = `/transactions?${params.toString()}`;
        });
    }

    // Reset filters
    if (resetFiltersBtn) {
        resetFiltersBtn.addEventListener('click', function() {
            window.location.href = '/transactions';
        });
    }

    // Set date inputs to current month range by default
    setDefaultDateRange();

    // Functions
    function openModal() {
        deleteModal.classList.add('active');
    }

    function closeModal() {
        deleteModal.classList.remove('active');
        transactionIdToDelete = null;
    }

    function setDefaultDateRange() {
        const startDateInput = document.getElementById('start-date');
        const endDateInput = document.getElementById('end-date');
        
        if (startDateInput && endDateInput) {
            // If no dates are set in the URL, set default dates
            const urlParams = new URLSearchParams(window.location.search);
            if (!urlParams.has('start_date') && !urlParams.has('end_date')) {
                const today = new Date();
                const firstDay = new Date(today.getFullYear(), today.getMonth(), 1);
                const lastDay = new Date(today.getFullYear(), today.getMonth() + 1, 0);
                
                startDateInput.value = firstDay.toISOString().split('T')[0];
                endDateInput.value = lastDay.toISOString().split('T')[0];
            }
        }
    }

    // Apply color to transaction amounts based on type
    const amountElements = document.querySelectorAll('.transaction-table tr');
    amountElements.forEach(row => {
        if (row.classList.contains('income')) {
            const amountCell = row.querySelector('.amount');
            if (amountCell) {
                amountCell.style.color = 'var(--income-color)';
            }
        } else if (row.classList.contains('expense')) {
            const amountCell = row.querySelector('.amount');
            if (amountCell) {
                amountCell.style.color = 'var(--expense-color)';
            }
        }
    });
});