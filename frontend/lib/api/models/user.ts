interface User {
  id: string;
  name: string;
  surname: string;
  email: string;
  created_at: string;
  updated_at: string;
}

interface UserResumed {
  name: string;
  surname: string;
  email: string;
}

export { User, UserResumed };
