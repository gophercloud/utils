package testing

import (
	"github.com/gophercloud/utils/gnocchi/metric/v1/resources"
)

// ResourceListResult represents raw server response from a server to a list call.
const ResourceListResult = `[
    {
        "created_by_project_id": "3d40ca37723449118987b9f288f4ae84",
        "created_by_user_id": "fdcfb420c09645e69e177a0bb1950884",
        "creator": "fdcfb420c09645e69e177a0bb1950884:3d40ca37723449118987b9f288f4ae84",
        "ended_at": null,
        "id": "1f3a0724-1807-4bd1-81f9-ee18c8ff6ccc",
        "metrics": {
            "cpu.delta": "2df1515e-6325-4d49-af0d-1052f6462fe4",
            "memory.usage": "777a01d6-4694-49cb-b86a-5ba9fd4e609e"
        },
        "original_resource_id": "1f3a0724-1807-4bd1-81f9-ee18c8ff6ccc",
        "project_id": "4154f08883334e0494c41155c33c0fc9",
        "revision_end": null,
        "revision_start": "2018-01-02T11:39:33.942419+00:00",
        "started_at": "2018-01-02T11:39:33.942391+00:00",
        "type": "compute_instance",
        "user_id": "bd5874d666624b24a9f01c128871e4ac"
    },
    {
        "created_by_project_id": "3d40ca37723449118987b9f288f4ae84",
        "created_by_user_id": "fdcfb420c09645e69e177a0bb1950884",
        "creator": "fdcfb420c09645e69e177a0bb1950884:3d40ca37723449118987b9f288f4ae84",
        "ended_at": null,
        "id": "789a7f65-977d-40f4-beed-f717100125f5",
        "metrics": {
            "disk.read.bytes.rate": "ed1bb76f-6ccc-4ad2-994c-dbb19ddccbae",
            "disk.write.bytes.rate": "0a2da84d-4753-43f5-a65f-0f8d44d2766c"
        },
        "original_resource_id": "789a7f65-977d-40f4-beed-f717100125f5",
        "project_id": "4154f08883334e0494c41155c33c0fc9",
        "revision_end": null,
        "revision_start": "2018-01-03T11:44:31.155773+00:00",
        "started_at": "2018-01-03T11:44:31.155732+00:00",
        "type": "compute_instance_disk",
        "user_id": "bd5874d666624b24a9f01c128871e4ac"
    }
]`

// Resource1 is an expected representation of a first resource from ResourceListResult.
var Resource1 = resources.Resource{
	CreatedByProjectID: "3d40ca37723449118987b9f288f4ae84",
	CreatedByUserID:    "fdcfb420c09645e69e177a0bb1950884",
	Creator:            "fdcfb420c09645e69e177a0bb1950884:3d40ca37723449118987b9f288f4ae84",
	ID:                 "1f3a0724-1807-4bd1-81f9-ee18c8ff6ccc",
	Metrics: map[string]string{
		"cpu.delta":    "2df1515e-6325-4d49-af0d-1052f6462fe4",
		"memory.usage": "777a01d6-4694-49cb-b86a-5ba9fd4e609e",
	},
	OriginalResourceID: "1f3a0724-1807-4bd1-81f9-ee18c8ff6ccc",
	ProjectID:          "4154f08883334e0494c41155c33c0fc9",
	RevisionStart:      "2018-01-02T11:39:33.942419+00:00",
	RevisionEnd:        "",
	StartedAt:          "2018-01-02T11:39:33.942391+00:00",
	EndedAt:            "",
	Type:               "compute_instance",
	UserID:             "bd5874d666624b24a9f01c128871e4ac",
}

// Resource2 is an expected representation of a second resource from ResourceListResult.
var Resource2 = resources.Resource{
	CreatedByProjectID: "3d40ca37723449118987b9f288f4ae84",
	CreatedByUserID:    "fdcfb420c09645e69e177a0bb1950884",
	Creator:            "fdcfb420c09645e69e177a0bb1950884:3d40ca37723449118987b9f288f4ae84",
	ID:                 "789a7f65-977d-40f4-beed-f717100125f5",
	Metrics: map[string]string{
		"disk.read.bytes.rate":  "ed1bb76f-6ccc-4ad2-994c-dbb19ddccbae",
		"disk.write.bytes.rate": "0a2da84d-4753-43f5-a65f-0f8d44d2766c",
	},
	OriginalResourceID: "789a7f65-977d-40f4-beed-f717100125f5",
	ProjectID:          "4154f08883334e0494c41155c33c0fc9",
	RevisionStart:      "2018-01-03T11:44:31.155773+00:00",
	RevisionEnd:        "",
	StartedAt:          "2018-01-03T11:44:31.155732+00:00",
	EndedAt:            "",
	Type:               "compute_instance_disk",
	UserID:             "bd5874d666624b24a9f01c128871e4ac",
}
