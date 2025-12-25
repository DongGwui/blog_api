# Domain Layer

클린 아키텍처의 핵심 계층으로, 비즈니스 로직의 본질을 정의합니다.

## 원칙

- **외부 의존성 금지**: 이 계층은 어떤 외부 패키지(프레임워크, DB 드라이버 등)에도 의존하지 않습니다
- **순수 Go 코드**: 표준 라이브러리만 사용
- **불변성 지향**: Entity는 가능한 불변으로 설계

## 디렉토리 구조

```
domain/
├── entity/           # 비즈니스 엔티티 (핵심 데이터 구조)
├── repository/       # Repository 인터페이스 (Port)
│   └── mocks/        # 테스트용 Mock 구현체
├── service/          # Service 인터페이스 (Port)
└── errors.go         # 도메인 에러 정의
```

## 파일별 역할

### entity/
비즈니스 도메인의 핵심 데이터 구조를 정의합니다.

| 파일 | 설명 |
|------|------|
| `category.go` | 카테고리 엔티티 |
| `tag.go` | 태그 엔티티 |
| `post.go` | 게시물 엔티티 (PostSummary, PostDetail 포함) |
| `project.go` | 프로젝트 엔티티 |
| `media.go` | 미디어 엔티티 (UploadedFile 포함) |
| `admin.go` | 관리자 엔티티 (TokenInfo, Claims 포함) |
| `dashboard.go` | 대시보드 통계 엔티티 |
| `view.go` | 조회수 결과 엔티티 |

### repository/
데이터 접근을 위한 인터페이스를 정의합니다. 구현체는 `infrastructure/` 계층에 있습니다.

| 파일 | 설명 |
|------|------|
| `category_repository.go` | 카테고리 CRUD 인터페이스 |
| `tag_repository.go` | 태그 CRUD 인터페이스 |
| `post_repository.go` | 게시물 CRUD + 검색 인터페이스 |
| `project_repository.go` | 프로젝트 CRUD + 정렬 인터페이스 |
| `media_repository.go` | 미디어 메타데이터 CRUD 인터페이스 |
| `storage_repository.go` | 파일 스토리지 인터페이스 (업로드/삭제/URL) |
| `admin_repository.go` | 관리자 계정 인터페이스 |
| `dashboard_repository.go` | 대시보드 통계 조회 인터페이스 |
| `view_repository.go` | 조회수 추적 인터페이스 (Redis용) |

### service/
비즈니스 로직을 위한 인터페이스를 정의합니다. 구현체는 `application/service/`에 있습니다.

| 파일 | 설명 |
|------|------|
| `category_service.go` | 카테고리 비즈니스 로직 인터페이스 |
| `tag_service.go` | 태그 비즈니스 로직 인터페이스 |
| `post_service.go` | 게시물 비즈니스 로직 인터페이스 |
| `project_service.go` | 프로젝트 비즈니스 로직 인터페이스 |
| `media_service.go` | 미디어 업로드/삭제 인터페이스 |
| `auth_service.go` | 인증 (로그인/JWT) 인터페이스 |
| `dashboard_service.go` | 대시보드 통계 인터페이스 |
| `view_service.go` | 조회수 기록 인터페이스 |

### errors.go
도메인 전용 에러를 정의합니다.

```go
// 사용 예시
if errors.Is(err, domain.ErrPostNotFound) {
    // 404 처리
}
```

## 새 도메인 추가 시

1. `entity/` 에 엔티티 정의
2. `repository/` 에 Repository 인터페이스 정의
3. `service/` 에 Service 인터페이스 정의
4. `errors.go` 에 도메인 에러 추가
5. `repository/mocks/` 에 Mock 구현체 추가 (테스트용)

## 코딩 컨벤션

```go
// Entity 예시
type Category struct {
    ID          int32
    Name        string
    Slug        string
    Description *string    // nullable 필드는 포인터
    PostCount   int64      // 집계 필드
    CreatedAt   time.Time
    UpdatedAt   *time.Time // nullable
}

// Repository 인터페이스 예시
type CategoryRepository interface {
    FindAll(ctx context.Context) ([]entity.Category, error)
    FindByID(ctx context.Context, id int32) (*entity.Category, error)
    Create(ctx context.Context, category *entity.Category) (*entity.Category, error)
    // ...
}

// Service 인터페이스 예시
type CategoryService interface {
    ListCategories(ctx context.Context) ([]entity.Category, error)
    CreateCategory(ctx context.Context, cmd *CreateCategoryCommand) (*entity.Category, error)
    // ...
}
```

## 테스트

- Mock은 `repository/mocks/`에 위치
- 인터페이스 기반으로 단위 테스트 작성 가능
- `application/service/` 테스트에서 Mock 사용
