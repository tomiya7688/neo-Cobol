#ifndef NEOC_RUNTIME_H
#define NEOC_RUNTIME_H

#ifdef __cplusplus
extern "C" {
#endif

/* Bootstrap runtime boundary. Shared runtime services will be added here
 * instead of leaking target-specific behavior into language semantics. */
void neoc_runtime_init(void);

#ifdef __cplusplus
}
#endif

#endif
