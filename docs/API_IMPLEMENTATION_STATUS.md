# Kasm Terraform Provider - API Implementation Status

This document tracks the implementation status of Kasm API endpoints in the Terraform provider.

## Status Legend
- = Implemented in provider
- = Implementation in progress
- = Not implemented

## Documented APIs

These APIs are officially documented in the Kasm API documentation.

### Resources

#### User Management
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/create_user | Implemented | kasm_user | internal/resources/user | ✅ | internal/resources/user/tests/user_test.go |
| POST /api/public/update_user | Implemented | kasm_user | internal/resources/user | ✅ | internal/resources/user/tests/user_test.go |
| DELETE /api/public/delete_user | Implemented | kasm_user | internal/resources/user | ✅ | internal/resources/user/tests/user_test.go |
| POST /api/public/update_user_attributes | Implemented | kasm_user | internal/resources/user | ✅ | internal/resources/user/tests/user_attributes_test.go |
| POST /api/public/logout_user | Implemented | kasm_user_logout | internal/resources/user | ✅ | internal/client/client_test.go |
| POST /api/public/get_user_attributes | Implemented | kasm_user | internal/resources/user | ✅ | internal/client/client_test.go |
| POST /api/public/import_user | Implemented | kasm_user | internal/resources/user | ✅ | internal/resources/user/tests/user_import_test.go |
| GET /api/public/get_users | Implemented | kasm_users | internal/datasources/users_list | ✅ | internal/datasources/users_list/tests/users_test.go |

#### Session Management
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/request_kasm | Implemented | kasm_session | internal/resources/session | ✅ | internal/resources/kasm/session/tests/session_test.go |
| POST /api/public/destroy_kasm | Implemented | kasm_session | internal/resources/session | ✅ | internal/resources/kasm/session/tests/session_test.go |
| POST /api/public/join_kasm | Implemented | kasm_join | internal/resources/join | ✅ | internal/resources/kasm/session/tests/session_test.go |
| POST /api/public/set_session_permissions | Implemented | kasm_session_permission | internal/resources/session_permission | ✅ | internal/resources/session_permission/tests/session_permission_basic_test.go |
| POST /api/public/keepalive | Implemented | kasm_keepalive | internal/resources/keepalive | ✅ | Unit: internal/resources/keepalive/resource_test.go, Acceptance: internal/resources/keepalive/tests/keepalive_test.go |
| POST /api/public/frame_stats | Not Implemented (Client Implementation Exists) | - | - | ❌ | - |
| POST /api/public/screenshot | Not Implemented (Client Implementation Exists) | - | - | ❌ | - |
| POST /api/public/exec_command | Not Implemented (Client Implementation Exists) | - | - | ❌ | - |
| POST /api/public/get_kasms | Implemented | kasm_sessions | internal/datasources/sessions | ✅ | internal/resources/kasm/session/tests/session_test.go |
| POST /api/public/get_kasm_status | Implemented | kasm_session_status | internal/datasources/session_status | ✅ | internal/resources/kasm/session/tests/session_test.go |
| GET /api/public/get_session_recordings | Not Implemented | - | - | ❌ | - |
| GET /api/public/get_sessions_recordings | Not Implemented | - | - | ❌ | - |
| POST /api/public/create_session | Implemented | kasm_session | internal/resources/session | ✅ | internal/resources/kasm/session/tests/session_test.go |

#### Image Management
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/create_image | Implemented | kasm_image | internal/resources/image | ✅ | internal/resources/image/tests/image_test.go |
| POST /api/public/update_image | Implemented | kasm_image | internal/resources/image | ✅ | internal/resources/image/tests/image_test.go |
| DELETE /api/public/delete_image | Implemented | kasm_image | internal/resources/image | ✅ | internal/resources/image/tests/image_test.go |
| GET /api/public/images | Implemented | kasm_images | internal/datasources/images | ✅ | internal/datasources/images/tests/images_test.go |

#### Registry Management
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/add_registry | Implemented | kasm_registry | internal/resources/registry | ✅ | internal/resources/registry/tests/registry_test.go |
| POST /api/public/remove_registry | Implemented | kasm_registry | internal/resources/registry | ✅ | internal/resources/registry/tests/registry_test.go |
| POST /api/public/get_registries | Implemented | kasm_registries | internal/datasources/registries | ✅ | internal/resources/registry/tests/registry_test.go |
| POST /api/public/update_registry | Implemented | kasm_registry | internal/resources/registry | ✅ | internal/resources/registry/tests/registry_test.go |

#### Cast Management
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/create_cast_config | Implemented | kasm_cast | internal/resources/cast | ✅ | internal/resources/cast/tests/cast_test.go |
| POST /api/public/update_cast_config | Implemented | kasm_cast | internal/resources/cast | ✅ | internal/resources/cast/tests/cast_test.go |
| DELETE /api/public/delete_cast_config | Implemented | kasm_cast | internal/resources/cast | ✅ | internal/resources/cast/tests/cast_test.go |

#### License Management
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/validate_license | Implemented | kasm_license | internal/client/license_ops.go | ✅ | internal/client/license_ops_test.go |
| POST /api/public/activate_license | Implemented | kasm_license | internal/client/license_ops.go | ✅ | internal/client/license_ops_test.go |
| POST /api/public/get_licenses | Implemented | kasm_license | internal/resources/license | ✅ | internal/client/license_ops_test.go |

#### Group Management
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/create_group | Implemented | kasm_group | internal/resources/group | ✅ | internal/resources/group/tests/group_test.go |
| POST /api/public/update_group | Implemented | kasm_group | internal/resources/group | ✅ | internal/resources/group/tests/group_test.go |
| DELETE /api/public/delete_group | Implemented | kasm_group | internal/resources/group | ✅ | internal/resources/group/tests/group_test.go |
| POST /api/public/set_group_membership | Implemented | kasm_group_membership | internal/resources/group_membership | ✅ | internal/resources/group_membership/tests/group_membership_test.go |

#### Group Image Management
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/create_group_image | Implemented | kasm_group_image | internal/resources/group_image | ✅ | internal/resources/group_image/tests/group_image_test.go |
| POST /api/public/update_group_image | Implemented | kasm_group_image | internal/resources/group_image | ✅ | internal/resources/group_image/tests/group_image_test.go |
| DELETE /api/public/delete_group_image | Implemented | kasm_group_image | internal/resources/group_image | ✅ | internal/resources/group_image/tests/group_image_test.go |
| POST /api/public/get_group_images | Implemented | kasm_group_images | internal/datasources/group_images | ✅ | internal/resources/group_image/tests/group_image_test.go |

#### RDP Client Connection
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/get_rdp_client_connection_info | Not Implemented | - | - | ❌ | - |

#### Egress Management
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/create_egress_provider | Not Implemented | - | - | ❌ | - |
| POST /api/public/update_egress_provider | Not Implemented | - | - | ❌ | - |
| DELETE /api/public/delete_egress_provider | Not Implemented | - | - | ❌ | - |
| POST /api/public/get_egress_providers | Not Implemented | - | - | ❌ | - |
| POST /api/public/create_egress_gateway | Not Implemented | - | - | ❌ | - |
| POST /api/public/update_egress_gateway | Not Implemented | - | - | ❌ | - |
| DELETE /api/public/delete_egress_gateway | Not Implemented | - | - | ❌ | - |
| POST /api/public/get_egress_gateways | Not Implemented | - | - | ❌ | - |
| POST /api/public/create_egress_credential | Not Implemented | - | - | ❌ | - |
| POST /api/public/update_egress_credential | Not Implemented | - | - | ❌ | - |
| DELETE /api/public/delete_egress_credential | Not Implemented | - | - | ❌ | - |
| POST /api/public/get_egress_credentials | Not Implemented | - | - | ❌ | - |
| POST /api/public/create_egress_provider_mapping | Not Implemented | - | - | ❌ | - |
| POST /api/public/update_egress_provider_mapping | Not Implemented | - | - | ❌ | - |
| DELETE /api/public/delete_egress_provider_mapping | Not Implemented | - | - | ❌ | - |
| POST /api/public/get_egress_provider_mappings | Not Implemented | - | - | ❌ | - |

#### Staging Configuration
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/create_staging_config | Not Implemented | - | - | ❌ | - |
| POST /api/public/update_staging_config | Not Implemented | - | - | ❌ | - |
| DELETE /api/public/delete_staging_config | Not Implemented | - | - | ❌ | - |
| POST /api/public/get_staging_configs | Not Implemented | - | - | ❌ | - |

### Data Sources

#### Images
| API Endpoint | Implementation Status | Data Source Name | File Location | Tests | Test File |
|--------------|---------------------|------------------|---------------|-------|-----------|
| GET /api/public/get_images | Implemented | kasm_images | internal/datasources/images | ✅ | internal/datasources/images/tests/images_test.go |

#### Registries
| API Endpoint | Implementation Status | Data Source Name | File Location | Tests | Test File |
|--------------|---------------------|------------------|---------------|-------|-----------|
| GET /api/public/get_registries | Implemented | kasm_registries | internal/datasources/registries | ✅ | internal/resources/registry/tests/registry_test.go |

#### Registry Images
| API Endpoint | Implementation Status | Data Source Name | File Location | Tests | Test File |
|--------------|---------------------|------------------|---------------|-------|-----------|
| GET /api/public/get_registry_images | Implemented | kasm_registry_images | internal/datasources/registry_images | ✅ | internal/datasources/registry_images/tests/registry_images_test.go |

#### Users
| API Endpoint | Implementation Status | Data Source Name | File Location | Tests | Test File |
|--------------|---------------------|------------------|---------------|-------|-----------|
| GET /api/public/get_user | Implemented | kasm_users | internal/datasources/users | ✅ | internal/resources/kasm/session/tests/session_test.go |
| GET /api/public/get_users | Implemented | kasm_users | internal/datasources/users_list | ✅ | internal/datasources/users_list/tests/users_test.go |

#### Zones
| API Endpoint | Implementation Status | Data Source Name | File Location | Tests | Test File |
|--------------|---------------------|------------------|---------------|-------|-----------|
| GET /api/public/get_zones | Implemented | kasm_zones | internal/datasources/zones | ❌ | - |

#### Licenses
| API Endpoint | Implementation Status | Data Source Name | File Location | Tests | Test File |
|--------------|---------------------|------------------|---------------|-------|-----------|
| GET /api/public/get_licenses | Implemented | kasm_licenses | internal/datasources/licenses | ❌ | - |

#### Groups
| API Endpoint | Implementation Status | Data Source Name | File Location | Tests | Test File |
|--------------|---------------------|------------------|---------------|-------|-----------|
| GET /api/public/get_groups | Implemented | kasm_groups | internal/datasources/groups | ✅ | internal/datasources/groups/tests/groups_test.go |
| GET /api/public/get_group | Implemented | kasm_groups | internal/datasources/groups | ❌ | - |

#### Group Memberships
| API Endpoint | Implementation Status | Data Source Name | File Location | Tests | Test File |
|--------------|---------------------|------------------|---------------|-------|-----------|
| POST /api/public/create_group_membership | Implemented | kasm_group_membership | internal/resources/group_membership | ✅ | internal/resources/group_membership/tests/group_membership_test.go |
| POST /api/public/update_group_membership | Implemented | kasm_group_membership | internal/resources/group_membership | ✅ | internal/resources/group_membership/tests/group_membership_test.go |
| DELETE /api/public/delete_group_membership | Implemented | kasm_group_membership | internal/resources/group_membership | ✅ | internal/resources/group_membership/tests/group_membership_test.go |
| POST /api/public/get_group_memberships | Implemented | kasm_group_memberships | internal/datasources/group_memberships | ✅ | internal/resources/group_membership/tests/group_membership_test.go |


## Undocumented APIs

These APIs are not officially documented in the Kasm API documentation but are available and used in the provider.

### Resources

#### Group Management
| API Endpoint | Implementation Status | Resource Name | File Location | Tests | Test File |
|--------------|---------------------|---------------|---------------|-------|-----------|
| POST /api/public/create_group | Implemented | kasm_group | internal/resources/group | ✅ | internal/resources/group/tests/group_test.go |
| POST /api/public/update_group | Implemented | kasm_group | internal/resources/group | ✅ | internal/resources/group/tests/group_test.go |
| DELETE /api/public/delete_group | Implemented | kasm_group | internal/resources/group | ✅ | internal/resources/group/tests/group_test.go |
| POST /api/public/set_group_membership | Implemented | kasm_group_membership | internal/resources/group_membership | ✅ | internal/resources/group_membership/tests/group_membership_test.go |

### Data Sources

#### Workspace
| API Endpoint | Implementation Status | Data Source Name | File Location | Tests | Test File |
|--------------|---------------------|------------------|---------------|-------|-----------|
| GET /api/public/get_workspace | In Progress | kasm_workspace | internal/datasources/workspace | ❌ | - |

## Missing Features

### Missing Data Sources (Documented APIs)
1. Sessions:
   - Need to create data source for `get_kasms` (client implementation exists)
   - Need to create data source for `get_kasm_status` (client implementation exists)
   - Need to create data source for `get_session_recordings` (client implementation exists)
   - Need to create data source for `get_sessions_recordings` (client implementation exists)

### Missing Resources (Documented APIs)
1. Session Features:
   - POST /api/public/frame_stats (for kasm_stats) - Client implementation exists
   - POST /api/public/screenshot (for kasm_screenshot) - Client implementation exists
   - POST /api/public/exec_command (for kasm_exec) - Client implementation exists
   - POST /api/public/get_rdp_client_connection_info (for kasm_rdp) - Client implementation exists

### Additional Undocumented Resources Found
1. Login Management:
   - Resource exists in internal/resources/login
2. Staging Management:
   - Resource exists in internal/resources/staging
3. Group Image Management:
   - Resource exists in internal/resources/group_image

## TODO
- [ ] Add version compatibility checks for undocumented APIs to ensure they work with different Kasm versions
- [ ] Consider implementing version-specific code paths for undocumented APIs if they change between versions
- [ ] Create data sources for session recordings functionality
- [ ] Add acceptance tests for session recordings (currently skipped in tests)

### User Import API

| Endpoint | Implemented in Code? | Tests Available? | Test File Location | Notes |
|----------|----------------------|------------------|--------------------|-------|
| User Import | Yes | Yes | [internal/resources/user/tests/user_import_test.go](cci:7://file:///Users/simon.garcia@contino.io/SynologyDrive/Code/HomeLab/GitHub/terraform-provider-kasm/internal/resources/user/tests/user_import_test.go:0:0-0:0) | Tests basic user import functionality with attribute verification |
