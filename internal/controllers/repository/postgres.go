package repository

import (
	"AvitoInternship/config"
	"AvitoInternship/internal/entity"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"time"
)

type PostgresDB struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgresDB(cfg *config.Config, l *zap.Logger) *PostgresDB {
	p, err := pgxpool.New(context.Background(), cfg.Postgres.Url)

	if err != nil {
		l.Error("Error received : ", zap.Error(err))
	}
	return &PostgresDB{p, l}
}

func (p *PostgresDB) CreateUser() (int, error) {
	q := "insert into users(id) values(default) RETURNING id"
	var id int
	err := p.pool.QueryRow(context.Background(), q).Scan(&id)

	if err != nil {
		e := "Received error while create new user in DB"
		p.logger.Error(e, zap.Error(err))
		return 0, err
	}
	return id, nil
}
func (p *PostgresDB) DeleteUser(id int) error {
	q := "delete from users where id = $1"
	t, err := p.pool.Exec(context.Background(), q, id)
	if err != nil {
		return err
	}
	if t.RowsAffected() < 1 {
		return UserNotFound{id}
	}
	return nil
}

func (p *PostgresDB) CreateSegment(segment *entity.Segment) error {
	q := "insert into segments values (DEFAULT, $1)"
	_, err := p.pool.Exec(context.Background(), q, segment.Slug)
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return SegmentAlreadyExist{segment.Slug}
		}
	}
	return err
}
func (p *PostgresDB) GetSegmentId(slug string) (int, error) {
	q := "select id from segments where slug = $1"
	var id int
	err := p.pool.QueryRow(context.Background(), q, slug).Scan(&id)
	if err != nil {
		e := fmt.Sprintf("Received error during search id for %s segment ", slug)
		p.logger.Error(e, zap.Error(err))
		return 0, SegmentNotFound{slug: slug}
	}
	return id, nil
}
func (p *PostgresDB) DeleteSegment(slug string) error {
	q := "delete from segments where slug = $1"
	t, err := p.pool.Exec(context.Background(), q, slug)
	if err != nil {
		return err
	}
	if t.RowsAffected() < 1 {
		return SegmentNotFound{slug}
	}
	return nil
}

func (p *PostgresDB) BindSegment(uid int, segment entity.Segment) error {
	var err error
	var q string
	if segment.EndDate == (time.Time{}) {
		q = "insert into user_segments values ($1,$2,$3)"
		_, err = p.pool.Exec(context.Background(), q, uid, segment.Id, segment.BeginDate)
	} else {
		q = "insert into user_segments values ($1,$2,$3,$4)"
		_, err = p.pool.Exec(context.Background(), q, uid, segment.Id, segment.BeginDate, segment.EndDate)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return SegmentAlreadyExist{segment.Slug}
		}
	}
	if err != nil {
		e := fmt.Sprintf("Received error while bind segment %s to user with id=%d", segment.Slug, uid)
		p.logger.Error(e, zap.Error(err), zap.Any("user", uid))
		return err
	}
	return nil
}
func (p *PostgresDB) UnBindSegment(uid int, segment entity.Segment) error {
	q := "update user_segments set end_date=$1 where uid=$2 and sid=$3 and begin_date=$4"
	_, err := p.pool.Exec(context.Background(), q, segment.EndDate, uid, segment.Id, segment.BeginDate)
	if err != nil {
		e := fmt.Sprintf("Received error while unbind segment %s from user with id=%d", segment.Slug, uid)
		p.logger.Error(e, zap.Error(err), zap.Any("user", uid))
		return err
	}
	return nil
}
func (p *PostgresDB) GetUserSegments(id int) ([]entity.Segment, error) {
	q := "select slug, begin_date, end_date from user_segments, segments where sid = id and uid = $1"
	rows, err := p.pool.Query(context.Background(), q, id)
	if err != nil {
		return nil, err
	}
	var segs []entity.Segment
	segs, err = pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Segment])
	if err != nil {
		return nil, err
	}
	return segs, nil
}
func (p *PostgresDB) GetUserActiveSegments(id int) ([]entity.Segment, error) {
	q := "select id, slug, begin_date, end_date from user_segments, segments where sid = id and uid = $1 and (end_date > $2 or end_date is null)"
	rows, err := p.pool.Query(context.Background(), q, id, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		e := fmt.Sprintf("Received error while get active segments for user with id=%d", id)
		p.logger.Error(e, zap.Error(err), zap.Any("user", id))
		return nil, err
	}
	var dtoSegs []PgSegmentDTO
	dtoSegs, err = pgx.CollectRows(rows, pgx.RowToStructByPos[PgSegmentDTO])
	if err != nil {
		e := fmt.Sprintf("Received error while unmarshall active segments for user with id=%d", id)
		p.logger.Error(e, zap.Error(err), zap.Any("user", id))
		return nil, err
	}
	segs := make([]entity.Segment, len(dtoSegs))
	for i := range segs {
		segs[i] = p.getSegment(&dtoSegs[i])
	}
	return segs, nil
}
func (p *PostgresDB) GetUserInactiveSegments(id int) ([]entity.Segment, error) {
	q := "select slug, begin_date, end_date from user_segments, segments where sid = id and uid = $1 and end_date <= $2"
	rows, err := p.pool.Query(context.Background(), q, id, time.Now().Format("2006-01-02 15:04:05"))
	if err != nil {
		e := fmt.Sprintf("Received error while get inactive segments for user with id=%d", id)
		p.logger.Error(e, zap.Error(err), zap.Any("user", id))
		return nil, err
	}
	var segs []entity.Segment
	segs, err = pgx.CollectRows(rows, pgx.RowToStructByPos[entity.Segment])
	if err != nil {
		e := fmt.Sprintf("Received error while unmarshall active segments for user with id=%d", id)
		p.logger.Error(e, zap.Error(err), zap.Any("user", id))
		return nil, err
	}
	return segs, nil
}
func (p *PostgresDB) getSegment(dto *PgSegmentDTO) entity.Segment {
	return entity.Segment{Id: dto.Id, Slug: dto.Slug, BeginDate: dto.BeginDate.Time, EndDate: dto.EndDate.Time}
}
func (p *PostgresDB) GetUserActiveSegment(uid int, slug string) (*entity.Segment, error) {
	q := " select segments.id, slug, begin_date, end_date  from user_segments, users, segments where users.id=$1 and uid=$1 and slug=$2 and sid = segments.id and (end_date > now() or end_date is null)"
	rows, err := p.pool.Query(context.Background(), q, uid, slug)
	if err != nil {
		e := fmt.Sprintf("Received error while get active %s segment for user with id=%d", slug, uid)
		p.logger.Error(e, zap.Error(err), zap.Any("user", uid), zap.Any("slug", slug))
		return nil, err
	}
	var dtoSeg PgSegmentDTO
	dtoSeg, err = pgx.CollectOneRow(rows, pgx.RowToStructByPos[PgSegmentDTO])
	if err != nil {
		e := fmt.Sprintf("Received error while unmarshall active %s segment for user with id=%d", slug, uid)
		p.logger.Error(e, zap.Error(err), zap.Any("user", uid), zap.Any("slug", slug))
		return nil, err
	}

	return &entity.Segment{Id: dtoSeg.Id, Slug: dtoSeg.Slug, BeginDate: dtoSeg.BeginDate.Time}, nil
}

func (p *PostgresDB) GetUserSegmentHistory(id int, sl string, af time.Time, bf time.Time) ([]entity.Segment, error) {
	f := "2006-01-02 15:04:05"
	q := fmt.Sprintf("%s%s%s",
		"select sid, slug, begin_date, end_date from user_segments, segments ",
		"where uid = $1 and slug = $2 and sid = segments.id and ",
		"(begin_date >= $3 or (end_date between $3 and $4) or end_date is null)")
	rows, err := p.pool.Query(context.Background(), q, id, sl, af.Format(f), bf.Format(f))
	if err != nil {
		e := fmt.Sprintf("Received error while get segments history for user")
		p.logger.Error(e, zap.Error(err), zap.Any("user", id), zap.Any("slug", sl))
		return nil, err
	}
	var dtoSegs []PgSegmentDTO
	dtoSegs, err = pgx.CollectRows(rows, pgx.RowToStructByPos[PgSegmentDTO])
	if err != nil {
		e := fmt.Sprintf("Received error while unmarshall active segments for user with id=%d", id)
		p.logger.Error(e, zap.Error(err), zap.Any("user", id))
		return nil, err
	}
	segs := make([]entity.Segment, len(dtoSegs))
	for i := range segs {
		segs[i] = p.getSegment(&dtoSegs[i])
	}
	return segs, nil
}
