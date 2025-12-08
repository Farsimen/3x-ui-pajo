/**
 * RBAC (Role-Based Access Control) - UI Control
 * Hides/shows menu items and features based on user role
 */

(function() {
    'use strict';

    // Cache for user role
    let userRole = null;
    let isAdmin = false;
    let isVendor = false;
    let vendorInbounds = [];

    /**
     * Fetch current user's role from server
     */
    function fetchUserRole() {
        return $.ajax({
            url: basePath + 'getUserRole',
            type: 'POST',
            dataType: 'json',
            async: false, // Synchronous to ensure role is loaded before page renders
            success: function(response) {
                if (response && response.obj) {
                    userRole = response.obj.role;
                    isAdmin = response.obj.isAdmin || false;
                    isVendor = response.obj.isVendor || false;
                    vendorInbounds = response.obj.inbounds || [];
                    
                    console.log('[RBAC] User role:', userRole);
                    console.log('[RBAC] Is Admin:', isAdmin);
                    console.log('[RBAC] Is Vendor:', isVendor);
                    if (isVendor) {
                        console.log('[RBAC] Vendor Inbounds:', vendorInbounds);
                    }
                }
            },
            error: function(xhr) {
                console.error('[RBAC] Failed to fetch user role:', xhr);
                // Default to restricted access if role fetch fails
                userRole = 'user';
                isAdmin = false;
                isVendor = false;
            }
        });
    }

    /**
     * Hide menu items based on user role
     */
    function applyMenuRestrictions() {
        if (!userRole) return;

        // For vendors, hide admin-only menu items
        if (isVendor && !isAdmin) {
            console.log('[RBAC] Applying vendor restrictions...');
            
            // Hide these menu items for vendors:
            // - Settings
            // - Xray Settings
            // - Any other admin-specific menus
            
            // Hide by checking href or data attributes
            $('a[href*="/panel/settings"]').closest('li').hide();
            $('a[href*="/panel/xray"]').closest('li').hide();
            $('a[href*="/panel/vendors"]').closest('li').hide();
            
            // Optionally hide other admin features
            $('.admin-only').hide();
            
            console.log('[RBAC] Vendor restrictions applied');
        }
        
        // For non-admin users (including vendors), hide Vendors menu
        if (!isAdmin) {
            $('a[href*="/panel/vendors"]').closest('li').hide();
        }
        
        // Show vendor menu if admin
        if (isAdmin) {
            console.log('[RBAC] Admin detected - showing all menus');
        }
    }

    /**
     * Filter inbounds list for vendors
     * Only show inbounds they have access to
     */
    function filterInboundsForVendor() {
        if (!isVendor || isAdmin) return;
        
        console.log('[RBAC] Filtering inbounds for vendor...');
        
        // This will be called after inbounds are loaded
        // Filter table rows to only show allowed inbounds
        if (vendorInbounds.length > 0) {
            $('#inbound-table tbody tr').each(function() {
                const inboundId = $(this).data('inbound-id');
                if (inboundId && !vendorInbounds.includes(inboundId)) {
                    $(this).hide();
                }
            });
        }
    }

    /**
     * Add visual indicator for vendor role
     */
    function addVendorIndicator() {
        if (isVendor && !isAdmin) {
            // Add badge to navbar showing vendor status
            const vendorBadge = $('<span>').
                addClass('badge badge-info ml-2')
                .text('Vendor Mode');
            
            $('.navbar .navbar-nav').prepend(
                $('<li>').addClass('nav-item').append(
                    $('<a>').addClass('nav-link').append(vendorBadge)
                )
            );
            
            console.log('[RBAC] Vendor indicator added');
        }
    }

    /**
     * Initialize RBAC system
     */
    function initRBAC() {
        console.log('[RBAC] Initializing...');
        
        // Fetch user role
        fetchUserRole();
        
        // Apply restrictions
        applyMenuRestrictions();
        
        // Add visual indicators
        addVendorIndicator();
        
        // Store in window for access by other scripts
        window.RBAC = {
            userRole: userRole,
            isAdmin: isAdmin,
            isVendor: isVendor,
            vendorInbounds: vendorInbounds,
            filterInbounds: filterInboundsForVendor
        };
        
        console.log('[RBAC] Initialization complete');
    }

    // Initialize when DOM is ready
    $(document).ready(function() {
        initRBAC();
    });

    // Export for use in other scripts
    window.initRBAC = initRBAC;

})();
