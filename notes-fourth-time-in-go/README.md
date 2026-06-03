# Go vs. Node.js JWT Authentication: Reference Guide

This guide is designed for developers coming from **Node.js (Express)** to help them understand how authentication, token generation, parsing, and middleware are implemented in **Go (Chi)**.

---

## 1. Token Generation

In Node.js, payloads are dynamic objects. In Go, everything is statically typed.

### 🔄 Comparison
| Feature | Node.js (`jsonwebtoken`) | Go (`golang-jwt/jwt/v5`) |
| :--- | :--- | :--- |
| **Payload Definition** | Dynamic JavaScript object. | Statically typed `struct` implementing standard claims. |
| **Signing Key Type** | `string` | `[]byte` (byte slice) |

### 💻 Code Example

#### Node.js (Express)
```javascript
const jwt = require('jsonwebtoken');

const token = jwt.sign(
  { userId: "123", email: "user@example.com" }, 
  process.env.JWT_SECRET, 
  { expiresIn: '15m' }
);
```

#### Go (Chi)
In Go, we define a custom struct that embeds `jwt.RegisteredClaims` and use json tags to structure the payload:

```go
type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    jwt.RegisteredClaims
}

claims := Claims{
    UserID: userid,
    Email:  email,
    RegisteredClaims: jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
    },
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
tokenString, err := token.SignedString([]byte(jwtSecret))
```

---

## 2. Verifying & Parsing Tokens

Node.js handles algorithm verification automatically. In Go, you must supply a custom **Keyfunc** callback to verify the signing algorithm before parsing the token.

### 💡 The Analogy
* **Node.js** is like showing your ID to a bouncer. The bouncer automatically checks if it's fake or uses a different format.
* **Go** is like presenting your document to a clerk. You must first supply the clerk with a checklist (the **Keyfunc**) stating exactly how to check the stamp (e.g. *"Ensure it is signed using HMAC"*). Only then does the clerk check the signature and hand you the verified data.

### 💻 Code Example

#### Node.js (Express)
```javascript
try {
  const decoded = jwt.verify(token, process.env.JWT_SECRET);
  console.log(decoded.userId);
} catch (err) {
  // Handles token expired, invalid signature, etc.
}
```

#### Go (Chi)
In Go, we parse using a callback to validate the algorithm type:

```go
token, err := jwt.ParseWithClaims(
    tokenString,
    &Claims{},
    func(t *jwt.Token) (any, error) {
        // 1. Explicitly check that the algorithm is HMAC
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, ErrTokenInvalid
        }
        // 2. Return the secret key as a byte slice
        return []byte(jwtSecret), nil
    },
)
```

---

## 3. Auth Middleware & Request Context

This is the biggest architectural shift. Node.js `req` objects are mutable. Go `http.Request` objects are immutable.

### 🔄 Comparison
| Feature | Node.js (Express) | Go (Standard `net/http`) |
| :--- | :--- | :--- |
| **Request Mutability** | **Mutable**. You can attach fields directly (`req.user = user`). | **Immutable**. You cannot modify the request object directly. |
| **Context Propagation** | Directly attaches properties to `req`. | Uses a typed, thread-safe context (`context.Context`). |
| **Middleware Chain** | Standard `next()` callback. | Wrap handler functions (`next.ServeHTTP(w, r.WithContext(ctx))`). |

### 💻 Code Example

#### Node.js (Express)
```javascript
const requireAuth = (req, res, next) => {
  const token = req.headers.authorization?.split(' ')[1];
  try {
    const decoded = jwt.verify(token, process.env.JWT_SECRET);
    req.userId = decoded.userId; // Directly mutate req
    next();
  } catch (err) {
    res.status(401).json({ error: 'Unauthorized' });
  }
};
```

#### Go (Chi)
Since Go's `http.Request` is immutable, we create a **new request context** with key-value pairs and pass it forward:

```go
// 1. Define custom, unexported key types to prevent namespace collisions
type contextKey string
const UserIDKey contextKey = "UserIDKey"

func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            token := extractToken(r)
            claims, err := utils.ParseToken(token, jwtSecret)
            if err != nil {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }

            // 2. Create a new context containing the claims data
            ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)

            // 3. Inject the context into a copy of the request and call next
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
```

#### Reading context values in Handlers:
```go
// Node.js
const userId = req.userId;

// Go
userId, ok := r.Context().Value(UserIDKey).(string) // Requires type assertion
```
