package authoring

import "strings"

func validateDemuxPatchCoverage(proposal DemuxProposal) error {
	sourceByID := map[string]HunkRange{}
	for _, hunk := range proposal.Hunks {
		id := strings.TrimSpace(hunk.ID)
		if id == "" {
			continue
		}
		sourceByID[id] = hunk
	}

	usedHunks := map[string]string{}
	wholeFileOwners := map[string]string{}
	for _, revision := range proposal.Revisions {
		revisionID := strings.TrimSpace(revision.ID)
		if revision.UseHunks {
			ids := revision.HunkIDs
			if len(ids) == 0 {
				ids = hunkIDs(revision.Hunks)
			}
			for _, id := range ids {
				id = strings.TrimSpace(id)
				if id == "" {
					continue
				}
				if owner, ok := usedHunks[id]; ok {
					return DuplicateHunkAssignmentError{HunkID: id, FirstOwnerID: owner, SecondOwnerID: revisionID}
				}
				source, ok := sourceByID[id]
				if !ok {
					return UnknownHunkReferenceError{RevisionID: revisionID, HunkID: id}
				}
				if referenced := findRevisionHunkByID(revision.Hunks, id); referenced.ID != "" && strings.TrimSpace(referenced.Patch) != "" && referenced.Patch != source.Patch {
					return HunkPatchBodyMismatchError{RevisionID: revisionID, HunkID: id, File: source.File}
				}
				usedHunks[id] = revisionID
			}
			continue
		}
		for _, file := range revision.Files {
			file = strings.TrimSpace(file)
			if file == "" {
				continue
			}
			if owner, ok := wholeFileOwners[file]; ok {
				return DuplicateWholeFileOwnerError{File: file, FirstOwnerID: owner, SecondOwnerID: revisionID}
			}
			wholeFileOwners[file] = revisionID
		}
	}
	return reviewDemuxHunkCoverage(proposal.Hunks, usedHunks, wholeFileOwners)
}
