# Application Layer

비즈니스 유스케이스를 구현하는 계층입니다. Domain 계층의 인터페이스를 구현합니다.

## 원칙

- **Domain 인터페이스 구현**: `domain/service/` 인터페이스를 구현
- **Repository 의존**: `domain/repository/` 인터페이스에 의존 (구현체가 아닌 인터페이스)
- **비즈니스 로직 집중**: 유효성 검사, 비즈니스 규칙, 트랜잭션 조율
- **외부 기술 중립**: HTTP, DB 세부사항을 알지 못함

## 디렉토리 구조

```
application/
└── service/          # Service 구현체
    ├── *_service.go       # 서비스 구현
    └── *_service_test.go  # 단위 테스트
```

## 파일별 역할

| 파일 | 구현 인터페이스 | 의존 Repository |
|------|----------------|-----------------|
| `category_service.go` | `CategoryService` | `CategoryRepository` |
| `tag_service.go` | `TagService` | `TagRepository` |
| `post_service.go` | `PostService` | `PostRepository` |
| `project_service.go` | `ProjectService` | `ProjectRepository` |
| `media_service.go` | `MediaService` | `MediaRepository`, `StorageRepository` |
| `auth_service.go` | `AuthService` | `AdminRepository` |
| `dashboard_service.go` | `DashboardService` | `DashboardRepository` |
| `view_service.go` | `ViewService` | `ViewRepository`, `ViewCountUpdater` |

## 서비스 구현 패턴

### 기본 구조

```go
package service

import (
    "context"
    "github.com/ydonggwui/blog-api/internal/domain"
    "github.com/ydonggwui/blog-api/internal/domain/entity"
    "github.com/ydonggwui/blog-api/internal/domain/repository"
    domainService "github.com/ydonggwui/blog-api/internal/domain/service"
)

type categoryService struct {
    repo repository.CategoryRepository
}

// 생성자: 인터페이스 반환
func NewCategoryService(repo repository.CategoryRepository) domainService.CategoryService {
    return &categoryService{repo: repo}
}

// 메서드 구현
func (s *categoryService) ListCategories(ctx context.Context) ([]entity.Category, error) {
    return s.repo.FindAll(ctx)
}
```

### 비즈니스 규칙 적용 예시

```go
func (s *categoryService) DeleteCategory(ctx context.Context, id int32) error {
    // 1. 존재 확인
    _, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return err  // 도메인 에러 또는 래핑된 에러 그대로 반환
    }

    // 2. 비즈니스 규칙: 게시물이 있으면 삭제 불가
    count, err := s.repo.GetPostCount(ctx, id)
    if err != nil {
        return fmt.Errorf("categoryService.DeleteCategory: %w", err)  // 에러 래핑
    }
    if count > 0 {
        return domain.ErrCategoryHasPosts
    }

    // 3. 삭제 실행
    if err := s.repo.Delete(ctx, id); err != nil {
        return fmt.Errorf("categoryService.DeleteCategory: %w", err)  // 에러 래핑
    }
    return nil
}
```

### Command 패턴 (생성/수정용)

```go
// domain/service/category_service.go 에 정의
type CreateCategoryCommand struct {
    Name        string
    Slug        string
    Description *string
}

// application/service/category_service.go 에서 사용
func (s *categoryService) CreateCategory(ctx context.Context, cmd *domainService.CreateCategoryCommand) (*entity.Category, error) {
    // 유효성 검사
    if exists, _ := s.repo.SlugExists(ctx, cmd.Slug); exists {
        return nil, domain.ErrCategorySlugExists
    }

    category := &entity.Category{
        Name:        cmd.Name,
        Slug:        cmd.Slug,
        Description: cmd.Description,
    }

    return s.repo.Create(ctx, category)
}
```

## 테스트 작성

Mock을 사용한 단위 테스트:

```go
func TestCategoryService_CreateCategory(t *testing.T) {
    // Mock 설정
    mockRepo := mocks.NewMockCategoryRepository()
    mockRepo.SlugExistsFunc = func(ctx context.Context, slug string) (bool, error) {
        return false, nil
    }
    mockRepo.CreateFunc = func(ctx context.Context, cat *entity.Category) (*entity.Category, error) {
        cat.ID = 1
        return cat, nil
    }

    // 서비스 생성
    svc := NewCategoryService(mockRepo)

    // 테스트 실행
    cmd := &domainService.CreateCategoryCommand{
        Name: "Test",
        Slug: "test",
    }
    result, err := svc.CreateCategory(context.Background(), cmd)

    // 검증
    assert.NoError(t, err)
    assert.Equal(t, int32(1), result.ID)
}
```

## 새 서비스 추가 시

1. `domain/service/` 에 인터페이스 정의
2. `application/service/` 에 구현체 작성
3. 생성자에서 Repository 인터페이스 주입
4. `domain/repository/mocks/` 에 Mock 추가
5. 단위 테스트 작성

## 에러 처리 패턴

### 에러 래핑 규칙

서비스 계층에서는 `fmt.Errorf("serviceName.FunctionName: %w", err)` 패턴으로 에러를 래핑합니다:

```go
func (s *postService) CreatePost(ctx context.Context, cmd *domainService.CreatePostCommand) (*entity.PostWithDetails, error) {
    // 도메인 에러는 그대로 반환
    if cmd.CategoryID != nil {
        _, err := s.categoryRepo.FindByID(ctx, *cmd.CategoryID)
        if err != nil {
            return nil, err  // domain.ErrCategoryNotFound 등 그대로 전달
        }
    }

    // 일반 에러는 래핑하여 반환
    post, err := s.repo.Create(ctx, entity)
    if err != nil {
        return nil, fmt.Errorf("postService.CreatePost: %w", err)
    }
    return post, nil
}
```

### 테스트에서 에러 비교

래핑된 에러를 비교할 때는 `errors.Is()` 사용:

```go
// 잘못된 방법
if err != domain.ErrPostNotFound { ... }

// 올바른 방법
if !errors.Is(err, domain.ErrPostNotFound) { ... }
```

## 주의사항

- HTTP 관련 코드 (gin.Context 등) 사용 금지
- DB 관련 코드 (sql.DB, sqlc 등) 사용 금지
- 도메인 에러 (`domain.ErrXxx`) 반환
- **일반 에러는 `fmt.Errorf`로 래핑** (함수명 포함)
- 에러 비교는 `errors.Is()` 사용
- context.Context는 항상 첫 번째 인자
