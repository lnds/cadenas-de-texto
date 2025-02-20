#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <stdbool.h>
#include <assert.h>
#include <time.h>

const int TEST_LIMIT = 10000000;
const int CHAR_MAX = 32; // lower case alphabet
const char MIN_CHAR = 'a';

bool naive_match(const char* pattern, const char* text);
void prepare_bitap(const char* pattern);
bool bitap_match(const char* pattern, const char *text);
bool bitap(const char* pattern, const char* text);
int* prepare_kmp(const char* pattern, int m);
bool kmp(int* lps, const char* pattern, int m, const char* text, int n);
void free_kmp(int *lps);
bool kmp_match(const char* pattern, const char* text);


int main(int argc, char* argv[]) {
  assert(naive_match("ab", "abacab"));
  assert(naive_match("aca", "abacab"));
  assert(naive_match("dabra", "abadacadabra"));
  assert(!naive_match("dabra", "abracadabera"));

  assert(bitap_match("ab", "abacab"));
  assert(bitap_match("aca", "abacab"));
  assert(bitap_match("dabra", "abadacadabra"));
  assert(!bitap_match("dabra", "abracadabera"));

  assert(kmp_match("ab", "abacab"));
  assert(kmp_match("aca", "abacab"));
  assert(kmp_match("dabra", "abadacadabra"));
  assert(!kmp_match("dabra", "abracadabera"));

  const char* pattern = "fragil";
  const char* text = "sdsajdlaksjdalsdjasldajsdlasjdlasdjaslkdasldsupercalifragilisticoespialidosojdkjadsdjslasjdlksjdasjd";

  clock_t elapsed = clock();
  for (int i = 0; i < TEST_LIMIT; i++) {
    assert(naive_match(pattern, text));
  }
  elapsed = clock() - elapsed;
  printf("elapsed time for naive: %f secs\n", (float)elapsed/CLOCKS_PER_SEC);

  prepare_bitap(pattern);
  elapsed = clock();
  for (int i = 0; i < TEST_LIMIT; i++) {
    assert(bitap(pattern, text));
  }
  elapsed = clock() - elapsed;
  printf("elapsed time for bitap: %f secs\n", (float)elapsed/CLOCKS_PER_SEC);
  
  int n = strlen(text);
  int m = strlen(pattern);
  int* lps = prepare_kmp(pattern, m);
  elapsed = clock();
  for (int i = 0; i < TEST_LIMIT; i++) {
    assert(kmp(lps, pattern, m, text, n));
  }
  elapsed = clock() - elapsed;
  printf("elapsed time for kmp: %f secs\n", (float)elapsed/CLOCKS_PER_SEC);
  free_kmp(lps);

  return 0;
}

bool naive_match(const char* pattern, const char* text) {
  int m = strlen(pattern);
  int n = strlen(text);
  for (int i = 0; i <= n-m; i++) {
    int j;
    for (j = 0; j < m; j++) {
      if (text[i+j] != pattern[j]) {
        break;
      }
    }
    if (j == m) {
      return true;
    }
  }
  return false;
}

unsigned long pattern_mask[CHAR_MAX+1];

void prepare_bitap(const char* pattern) {
  int m = strlen(pattern);
  assert(m < 32);
  for (int i = 0; i <= CHAR_MAX; i++) {
    pattern_mask[i] = ~0;
  }
  for (int i = 0; i < m; i++) {
    pattern_mask[pattern[i]-MIN_CHAR] &= ~(1UL << i);
  }
}

bool bitap(const char* pattern, const char* text) {
  int m = strlen(pattern);
  assert(m < 32);
  unsigned long bmask = (1UL << m);
  unsigned long r = ~1;
  for (int i = 0; text[i] != '\0'; ++i) {
    r |= pattern_mask[text[i]-MIN_CHAR];
    r <<= 1;
    if (0 == (r & bmask)) {
      return true;
    }
  }
  return false;
}

bool bitap_match(const char* pattern, const char* text) {
  prepare_bitap(pattern);
  return bitap(pattern, text);
}


int* prepare_kmp(const char* pattern, int m) {
  int *lps = (int*) malloc(m * sizeof(int));
  int pos = 1;
  int cnd = 0;
  lps[0] = -1;
  while (pos < m) {
    if (pattern[pos] == pattern[cnd]) {
      lps[pos] = cnd;
    } else {
      lps[pos] = cnd;
      while (cnd >= 0 && pattern[pos] != pattern[cnd]) {
        cnd = lps[cnd];
      }
    }
    pos++;
    cnd++;
  }
  return lps;
}

// kmp modified to return only first occurrence
bool kmp(int* lps, const char* pattern, int m, const char* text, int n) {
  int j = 0;
  int k = 0;
  while (j < n) {
    if (pattern[k] == text[j]) {
      j++;
      k++;
      if (k == m) {
        return true;
      }
    } else {
      k = lps[k];
      if (k < 0) {
        j++;
        k++;
      }
    }
  }
  return false;
}

void free_kmp(int *lps) {
  free(lps);
}

bool kmp_match(const char* pattern, const char* text) {
  int n = strlen(text);
  int m = strlen(pattern);
  int *lps = prepare_kmp(pattern, m);
  bool result = kmp(lps, pattern, m, text, n);
  free_kmp(lps);
  return result;
}
