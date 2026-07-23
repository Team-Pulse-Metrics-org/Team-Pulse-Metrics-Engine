import { useEffect, useState } from "react";

interface User {
  id: string;
  email: string;
  role: string;
  first_name: string;
  last_name: string;
}

const API = import.meta.env.VITE_API_URL;

export default function AdminPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [email, setEmail] = useState("");
  const [role, setRole] = useState("developer");
  // Fetch all users
  useEffect(() => {
    fetch(`${API}/api/v1/admin/users`)
      .then((res) => res.json())
      .then((data) => setUsers(data))
      .catch((err) => console.error("Error fetching users:", err));
  }, []);
const handleRoleChange = async (id: string, newRole: string) => {
  try {
    const response = await fetch(
      `${API}/api/v1/admin/users/${id}/role`,
      {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          role: newRole,
        }),
      }
    );

    if (!response.ok) {
      throw new Error("Failed to update role");
    }

    setUsers(
      users.map((user) =>
        user.id === id
          ? { ...user, role: newRole }
          : user
      )
    );
  } catch (error) {
    console.error(error);
  }
};
const handleRemove = async (id: string) => {
  try {
    await fetch(
      `${API}/api/v1/admin/users/${id}`,
      {
        method: "DELETE",
      }
    );

    setUsers(users.filter((user) => user.id !== id));
  } catch (error) {
    console.error("Error deleting user:", error);
  }
};
 const handleAddUser = async () => {
  try {
    const response = await fetch(
      `${API}/api/v1/admin/users`,
      {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          first_name: firstName,
          last_name: lastName,
          email,
          role,
        }),
      }
    );

    if (!response.ok) {
      const errorText = await response.text();
      console.log(errorText);
      throw new Error("Failed to create user");
    }

    const newUser = await response.json();

    setUsers([...users, newUser]);

    setFirstName("");
    setLastName("");
    setEmail("");
    setRole("developer");
  } catch (error) {
    console.error("Error creating user:", error);
  }
};
  return (
    <div className="p-6">
      {/* Header */}
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-4xl font-bold">Admin Panel</h1>
<button
  onClick={handleAddUser}
  className="bg-blue-600 hover:bg-blue-700 px-4 py-2 rounded-lg"
>
  + Add User
</button>
        
      </div>
<div className="grid grid-cols-4 gap-4 mb-6">
  <input
    type="text"
    placeholder="First Name"
    value={firstName}
    onChange={(e) => setFirstName(e.target.value)}
    className="border rounded p-2 bg-gray-800"
  />

  <input
    type="text"
    placeholder="Last Name"
    value={lastName}
    onChange={(e) => setLastName(e.target.value)}
    className="border rounded p-2 bg-gray-800"
  />

  <input
    type="email"
    placeholder="Email"
    value={email}
    onChange={(e) => setEmail(e.target.value)}
    className="border rounded p-2 bg-gray-800"
  />

  <select
    value={role}
    onChange={(e) => setRole(e.target.value)}
    className="border rounded p-2 bg-gray-800"
  >
    <option value="developer">Developer</option>
    <option value="lead">Team Leader</option>
    <option value="administrator">Admin</option>
  </select>
</div>

      {/* Users Table */}
      <table className="w-full border border-gray-700">
        <thead>
          <tr className="bg-gray-800">
            <th className="p-3 text-left">User</th>
            <th className="p-3 text-left">Email</th>
            <th className="p-3 text-left">Role</th>
            <th className="p-3 text-left">Action</th>
          </tr>
        </thead>

        <tbody>
          {users.map((user) => (
            <tr key={user.id} className="border-t border-gray-700">
              {/* User Name */}
              <td className="p-3">
                {user.first_name} {user.last_name}
              </td>

              {/* Email */}
              <td className="p-3">{user.email}</td>

              {/* Role Dropdown */}
              <td className="p-3">
                <select
                  value={user.role}
                  onChange={(e) =>
                    handleRoleChange(user.id, e.target.value)
                  }
                  className="bg-gray-800 border border-gray-600 rounded px-3 py-1"
                >
                  <option value="developer">Developer</option>
                  <option value="lead">Team Leader</option>
                  <option value="administrator">Admin</option>
                </select>
              </td>

              {/* Remove Button */}
              <td className="p-3">
        <button
            onClick={() => handleRemove(user.id)}
            className="bg-red-500 hover:bg-red-600 px-3 py-1 rounded"
              >
            Remove
        </button>
        </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}