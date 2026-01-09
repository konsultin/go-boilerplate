package repository

import (
	"github.com/konsultin/project-goes-here/internal/svc-core/model"
	"github.com/konsultin/errk"
)

func (r *Repository) FindSessionByXid(xid string) (*model.AuthSession, error) {
	var m model.AuthSession
	err := r.sql.AuthSession.FindByXid.GetContext(r.ctx, &m, xid)
	if err != nil {
		return nil, errk.Trace(err)
	}
	return &m, nil
}

func (r *Repository) DeleteSessionById(id int64) error {
	_, err := r.sql.AuthSession.DeleteById.ExecContext(r.ctx, id)
	if err != nil {
		return errk.Trace(err)
	}
	return nil
}

func (r *Repository) InsertAuthSession(session *model.AuthSession) error {
	err := r.sql.AuthSession.Insert.GetContext(r.ctx, &session.Id, session)
	if err != nil {
		return errk.Trace(err)
	}
	return nil
}

// GetDownloadFileUrl returns the download URL for a given file path
// TODO: Implement based on your file storage (S3, Local, etc.)
func (r *Repository) GetDownloadFileUrl(path string) (string, error) {
	// Placeholder implementation - replace with actual file storage logic
	// Example for S3: return s3Client.GetPresignedURL(path)
	// Example for local: return fmt.Sprintf("http://localhost:8080/files/%s", path), nil
	return "https://example.com/" + path, nil
}
