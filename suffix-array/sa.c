#include <stdio.h>
#include <string.h>
#include <stdlib.h>

struct suffix {
  int index;
  char *suff;
};

int cmp(const void *pa, const void *pb) {
  char* sa = ((struct suffix *) pa)->suff;
  char* sb = ((struct suffix *) pb)->suff;
  return strcmp(sa, sb);
}

int* build_suffix_array(char *txt, int n) {
  struct suffix* suffixes = (struct suffix*) malloc(n*sizeof(struct suffix));

  for (int i = 0; i < n; i++) {
    suffixes[i].index = i;
    suffixes[i].suff = (txt+i);
  }
  qsort(suffixes, n, sizeof(struct suffix), cmp);

  int *suffix_arr = (int*) malloc(n*sizeof(int));
  for (int i = 0; i < n; i++) {
    suffix_arr[i] = suffixes[i].index;
  }
  free(suffixes);
  return suffix_arr;
}

int search(char* pattern, char* txt, int *suffix_arr, int n) {
  int m = strlen(pattern);
  int l = 0, r = n-1;
  while (l <= r) {
    int mid = l + (r-l)/2;
    int res = strncmp(pattern, txt+suffix_arr[mid], m);
    if (res == 0) {
      return suffix_arr[mid];
    }
    if (res < 0) r = mid - 1;
    else l = mid + 1;
  }
  return -1;
}

char* read_text(const char* filename) {
  FILE* f = fopen(filename, "rb");
  if (!f) {
    return NULL;
  }
  fseek(f, 0, SEEK_END);
  int len = ftell(f);
  char* buf = (char*) malloc(len);
  if (buf) {
    fseek(f, 0, SEEK_SET);
    fread(buf, 1, len, f);
  }
  fclose(f);
  return buf;
}

int main(int argc, char* argv[]) {

  if (argc != 3) {
    printf("usage: sa pattern file\n");
    return -1;
  }


  char* txt = read_text(argv[2]);
  if (txt == NULL) {
    printf("file not found\n");
    return -1;
  }

  char* pattern = argv[1];

  int n = strlen(txt);
  int* sa = build_suffix_array(txt, n);
  printf("pattern found at position: %d\n",search(pattern, txt, sa, n));
  return 0;
}
