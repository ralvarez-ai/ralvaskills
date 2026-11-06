# Code Design Refactor Reference

Examples moved out of [SKILL.md](SKILL.md).

## Primitive obsession: value-object examples (§8)

```python
@dataclass(frozen=True, slots=True)
class Email:
    value: str
    def __post_init__(self) -> None:
        if "@" not in self.value:
            raise ValueError(f"invalid email: {self.value!r}")
```

```go
type Email string

func NewEmail(s string) (Email, error) {
    if !strings.Contains(s, "@") {
        return "", fmt.Errorf("invalid email: %q", s)
    }
    return Email(s), nil
}
```
