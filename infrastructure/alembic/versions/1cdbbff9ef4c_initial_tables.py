"""Initial tables

Revision ID: 1cdbbff9ef4c
Revises:
Create Date: 2026-01-29 08:13:12.670857

"""
from typing import Sequence, Union

from alembic import op
import sqlalchemy as sa
from sqlalchemy.dialects import postgresql


# revision identifiers, used by Alembic.
revision: str = '1cdbbff9ef4c'
down_revision: Union[str, None] = None
branch_labels: Union[str, Sequence[str], None] = None
depends_on: Union[str, Sequence[str], None] = None




def upgrade() -> None:
    """Upgrade schema."""
    # Create the base schema
    op.execute('CREATE SCHEMA IF NOT EXISTS base')

    # Create enum types using raw SQL for proper IF NOT EXISTS handling
    op.execute("""
        DO $$ BEGIN
            CREATE TYPE base.cpu_architecture AS ENUM ('amd64', 'arm64');
        EXCEPTION
            WHEN duplicate_object THEN null;
        END $$;
    """)

    op.execute("""
        DO $$ BEGIN
            CREATE TYPE base.simulation_run_status AS ENUM ('pending', 'queued', 'running', 'completed', 'failed');
        EXCEPTION
            WHEN duplicate_object THEN null;
        END $$;
    """)

    # Create client table
    op.create_table(
        'client',
        sa.Column('id', sa.String, primary_key=True),
        sa.Column('created_at', sa.DateTime(timezone=True), nullable=False),
        sa.Column('last_active_at', sa.DateTime(timezone=True), nullable=False),
        schema='base'
    )

    # Create simulation_meta table
    op.create_table(
        'simulation_meta',
        sa.Column('id', sa.String, primary_key=True),
        sa.Column('name', sa.String, nullable=False),
        sa.Column('description', sa.Text, nullable=True),
        sa.Column('parameters', postgresql.JSONB, nullable=True),
        sa.Column('outputs', postgresql.JSONB, nullable=True),
        sa.Column('created_at', sa.DateTime(timezone=True), nullable=False),
        sa.Column('updated_at', sa.DateTime(timezone=True), nullable=False),
        schema='base',
    )

    # Create simulation_binary table
    op.create_table(
        'simulation_binary',
        sa.Column('simulation_id', sa.String, nullable=False),
        sa.Column('cpu_arch', postgresql.ENUM('amd64', 'arm64', name='cpu_architecture', schema='base', create_type=False), nullable=False),
        sa.Column('blob', sa.LargeBinary, nullable=True),
        sa.Column('revision', sa.Integer, nullable=False, server_default='0'),
        sa.Column('created_at', sa.DateTime(timezone=True), nullable=False),
        sa.Column('updated_at', sa.DateTime(timezone=True), nullable=False),
        sa.PrimaryKeyConstraint('simulation_id', 'cpu_arch'),
        sa.ForeignKeyConstraint(
            ['simulation_id'],
            ['base.simulation_meta.id'],
            ondelete='CASCADE'
        ),
        schema='base'
    )

    # Create simulation_run table
    op.create_table(
        'simulation_run',
        sa.Column('client_id', sa.String, primary_key=True),
        sa.Column('simulation_id', sa.String, nullable=True),
        sa.Column('status', postgresql.ENUM('pending', 'queued', 'running', 'completed', 'failed', name='simulation_run_status', schema='base', create_type=False), nullable=False, server_default='pending'),
        sa.Column('parameters', postgresql.JSONB, nullable=True),
        sa.Column('created_at', sa.DateTime(timezone=True), nullable=False),
        sa.Column('updated_at', sa.DateTime(timezone=True), nullable=False),
        sa.Column('completed_at', sa.DateTime(timezone=True), nullable=True),
        sa.ForeignKeyConstraint(
            ['client_id'],
            ['base.client.id'],
            ondelete='CASCADE'
        ),
        sa.ForeignKeyConstraint(
            ['simulation_id'],
            ['base.simulation_meta.id'],
            ondelete='SET NULL'
        ),
        schema='base'
    )

    # Create simulation_output table
    op.create_table(
        'simulation_output',
        sa.Column('client_id', sa.String, primary_key=True),
        sa.Column('blob', sa.LargeBinary, nullable=True),
        sa.Column('created_at', sa.DateTime(timezone=True), nullable=False),
        sa.Column('updated_at', sa.DateTime(timezone=True), nullable=False),
        sa.ForeignKeyConstraint(
            ['client_id'],
            ['base.client.id'],
            ondelete='CASCADE'
        ),
        schema='base'
    )

    # Create api_key table
    op.create_table(
        'api_key',
        sa.Column('id', sa.String, primary_key=True),
        sa.Column('owner', sa.String, nullable=False),
        sa.Column('key', sa.String, nullable=False, unique=True),
        sa.Column('revoked', sa.Boolean, nullable=False, server_default='false'),
        sa.Column('expires_at', sa.DateTime(timezone=True), nullable=False),
        sa.Column('created_at', sa.DateTime(timezone=True), nullable=False),
        schema='base'
    )

    # Create indexes for common queries
    op.create_index(
        'ix_simulation_run_status',
        'simulation_run',
        ['status'],
        schema='base'
    )
    op.create_index(
        'ix_api_key_key',
        'api_key',
        ['key'],
        schema='base'
    )
    op.create_index(
        'ix_client_last_active_at',
        'client',
        ['last_active_at'],
        schema='base'
    )


def downgrade() -> None:
    """Downgrade schema."""

    op.execute('DROP SCHEMA IF EXISTS base CASCADE')
